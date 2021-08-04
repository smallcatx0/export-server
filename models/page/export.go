package page

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"export-server/bootstrap/global"
	"export-server/models/dao"
	"export-server/models/dao/mdb"
	"export-server/models/dao/rdb"
	"export-server/pkg/conf"
	"export-server/pkg/exception"
	"export-server/pkg/glog"
	"export-server/pkg/helper"
	"export-server/valid"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"gorm.io/gorm"
)

type ExportServ struct{}

// Detail 导出详情查询
func (e *ExportServ) Detail(key string) (ret interface{}, err error) {
	explog := make(map[string]interface{}, 5)
	res := dao.MDB.
		Model(&mdb.ExportLog{}).
		Select([]string{"id", "hash_key", "title", "ext_type", "source_type", "status", "fail_reason", "created_at"}).
		First(&explog, "hash_key = ?", key)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err = exception.ExNotFund
		return
	}
	if !helper.EqualInt(explog["status"], mdb.ExportLog_status_succ) {
		ret = explog
		return
	}
	explog["down_url"] = new(mdb.ExportFile).DownUrl(key)
	ret = explog
	return
}

func (e *ExportServ) History(c *gin.Context, param *valid.ExpLogHistory) (ret interface{}, err error) {
	explogs := make([]map[string]interface{}, 0, 10)
	last7d := carbon.Now().SubDays(7).ToDateString()
	res := dao.MDB.Debug().
		Model(&mdb.ExportLog{}).
		Select([]string{"id", "hash_key", "title", "ext_type", "source_type", "status", "fail_reason", "created_at"}).
		Order("id DESC").
		Find(&explogs, "user_id = ? AND created_at > ?", param.Uid, last7d)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err = res.Error
		return
	}
	ret = explogs
	return
}

// 获取参数哈希
func (e *ExportServ) paramHash(v interface{}) (hashKey string, err error) {
	paramBt, err := json.Marshal(v)
	if err != nil {
		return
	}
	hash := md5.Sum(paramBt)
	hashKey = fmt.Sprintf("%x", hash)
	return
}

func (e *ExportServ) HandelSHttp(c *gin.Context, param *valid.ExpSHttpParam) (ret interface{}, err error) {
	// 1. 参数哈希
	hashKey, err := e.paramHash(param)
	if err != nil {
		glog.Error("param hash err", "", err.Error())
	}
	ret = map[string]string{"hash_key": hashKey}

	// 2. 查询任务是否存在，不存在则记录 存在直接返回
	if new(mdb.ExportLog).HashHeyExisted(hashKey) {
		glog.Info("任务已经存在 hash_key=" + hashKey)
		return
	}
	expLog := &mdb.ExportLog{
		HashKey:    hashKey,
		Title:      param.Title,
		ExtType:    param.EXTType,
		SourceType: mdb.ExportLog_Stype_Http,
		Status:     mdb.ExportLog_status_pending,
		Callback:   param.CallBack,
		UserId:     param.UserID,
	}
	source, err := json.Marshal(param.SourceHTTP)
	if err != nil {
		glog.Error("json.Marshal err", "", err.Error())
		return
	}
	expLog.Param = string(source)
	res := dao.MDB.Create(expLog)
	if res.Error != nil {
		glog.Error("exportLog insert err", "", res.Error.Error())
		err = res.Error
		return
	}

	// 3. 准备参数丢任务队列中
	httpQ := &rdb.Mq{
		Key: global.TaskHttpKey,
	}
	httpQ.Push(&rdb.ExportTask{
		TaskID: hashKey,
	})
	return
}

func (e *ExportServ) HandelSRaw(c *gin.Context, param *valid.ExpSRawParam) (ret interface{}, err error) {
	// 1. 参数哈希
	hashKey, err := e.paramHash(param)
	if err != nil {
		glog.Error("param hash err", "", err.Error())
		return
	}
	ret = map[string]string{"hash_key": hashKey}

	// 2. 查询任务是否存在，不存在则记录 存在直接返回
	if new(mdb.ExportLog).HashHeyExisted(hashKey) {
		glog.Info("任务已经存在 hash_key=" + hashKey)
		return
	}
	// 将参数中的source_raw存到本地文件中
	paramDir := conf.AppConf.GetString("storage.source_raw")
	paramSavePath := path.Join(paramDir, hashKey+".json")
	helper.TouchDir(paramSavePath)
	err = ioutil.WriteFile(paramSavePath, []byte(param.SourceRaw), 0666)
	if err != nil {
		glog.Error("writeFile err", "", err.Error())
		return
	}

	expLog := &mdb.ExportLog{
		HashKey:    hashKey,
		Title:      param.Title,
		ExtType:    param.EXTType,
		SourceType: mdb.ExportLog_Stype_Raw,
		Callback:   param.CallBack,
		Status:     mdb.ExportLog_status_pending,
		UserId:     param.UserID,
	}
	res := dao.MDB.Create(expLog)
	if res.Error != nil {
		glog.Error("exportLog insert err", "", res.Error.Error())
		err = res.Error
		return
	}

	// 3. 准备参数丢任务队列中
	httpQ := &rdb.Mq{
		Key: global.TaskRawKey,
	}
	httpQ.Push(&rdb.ExportTask{
		TaskID: hashKey,
	})
	return
}
