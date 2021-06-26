package page

import (
	"crypto/md5"
	"encoding/json"
	"export-server/models/dao"
	"export-server/models/dao/mdb"
	"export-server/pkg/glog"
	"export-server/valid"
	"strings"

	"github.com/gin-gonic/gin"
)

type ExportServ struct{}

func (e *ExportServ) Handel(c *gin.Context, param *valid.ExportParam) (data interface{}, err error) {
	// 1. 获取参数的哈希
	// TODO: 当数据源为直传时，要不要将 SourceRaw 也计算到哈希里去
	paramBt, err := json.Marshal(param)
	if err != nil {
		glog.Error("json.Marshal err", "", err.Error())
		return
	}
	hash := md5.Sum(paramBt)
	hashKey := string(hash[:])
	data = map[string]string{"hash_key": hashKey}
	// 2. 记录请求日志
	err = e.RecordLog(hashKey, param)
	if err != nil {
		return
	}
	// 3. 准备参数丢任务队列中
	switch strings.ToLower(param.SourceType) {
	case "http":
	}
	return
}

func (e *ExportServ) RecordLog(hashKey string, param *valid.ExportParam) error {
	// 存数据库

	expLog := &mdb.ExportLog{
		HashKey:    hashKey,
		Title:      param.Title,
		ExtType:    param.EXTType,
		SourceType: param.SourceType,
		Callback:   param.CallBack,
		UserId:     param.UserID,
	}
	switch strings.ToLower(param.EXTType) {
	case "http":
		sourse, err := json.Marshal(param.SourceHTTP)
		if err != nil {
			glog.Error("json.Marshal", "", err.Error())
			return err
		}
		expLog.Param = string(sourse)
	case "sql":
		sourse, err := json.Marshal(param.SourceSQL)
		if err != nil {
			glog.Error("json.Marshal", "", err.Error())
			return err
		}
		expLog.Param = string(sourse)
	}
	res := dao.MDB.Create(expLog)
	if res.Error != nil {
		glog.Error("exportlog insert err", "", res.Error.Error())
	}
	return nil
}
