package data

import (
	"export-server/bootstrap/global"
	"export-server/models/dao"
	cal "export-server/models/dao/Cal"
	"export-server/models/dao/aoss"
	"export-server/models/dao/mdb"
	"export-server/models/dao/rdb"
	"export-server/pkg/conf"
	"export-server/pkg/excel"
	"export-server/pkg/glog"
	"export-server/pkg/helper"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type RawWorker struct {
	Tasks  *rdb.Mq
	taskCh chan *rdb.ExportTask
}

func (w *RawWorker) Run(pool int) {
	// 单消费端 多任务执行
	w.Tasks = &rdb.Mq{Key: global.TaskRawKey}
	// 缓冲区越大，程序宕机后丢消息越多
	w.taskCh = make(chan *rdb.ExportTask, 20)
	log.Print("RawWorker pool=", pool)
	// 启动工作协程
	w.startWorker(pool)

	// 监听队列
	go func() {
		w.Tasks.BPop(func(s string) {
			atask := &rdb.ExportTask{}
			atask.Build(s)
			// 丢入缓冲区
			w.taskCh <- atask
		})
	}()
}

// startWorker 启动工作协程
func (w *RawWorker) startWorker(pool int) {
	for i := 0; i < pool; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					glog.Error("[runtime] err ,recoverd", "", fmt.Errorf("%v", err).Error())
				}
			}()
			for {
				w.work()
			}
		}()
	}
}

func (w *RawWorker) work() {
	currTask := <-w.taskCh
	taskID := currTask.TaskID
	// 1. 数据库中查询任务详情
	expLog := mdb.ExportLog{}
	result := dao.MDB.Where("hash_key=?", taskID).First(&expLog)
	if result.Error != nil {
		glog.Error("TaskNotFund hash_key=" + taskID)
		return
	}
	// 任务取消
	if expLog.Status == mdb.ExportLog_status_cancle {
		return
	}
	// 2. 拿到json数据 -> 3. 写入excel
	paramDir := conf.AppConf.GetString("tmp_storage.source_raw")
	excelTmpPath := conf.AppConf.GetString("tmp_storage.outexcel_tmp")
	filename := expLog.Title + "." + expLog.ExtType
	paramFilePath := path.Join(paramDir, taskID+".json")
	lists, err := ioutil.ReadFile(paramFilePath)
	if err != nil {
		reason := "请求参数json文件读取失败" + err.Error()
		expLog.SaveFailReason(reason)
		return
	}
	excelw := excel.NewExcelRecorder(path.Join(excelTmpPath, taskID, filename))
	excelw.JsonListWrite(excel.Pos{X: 1, Y: 1}, string(lists), true)
	excelw.Save()

	// 4. 压缩文件夹 并删除源文件
	zipFilePath := path.Join(excelTmpPath, taskID+".zip")
	taskDir := path.Join(excelTmpPath, taskID)
	helper.FolderZip(taskDir, zipFilePath)
	glog.ErrorOnly("remove Files err path="+taskDir, "", os.RemoveAll(taskDir))
	glog.ErrorOnly("remove Files err path="+paramFilePath, "", os.Remove(paramFilePath))

	// 5. 上传云 OOS -> 删除本地文件
	objname, err := aoss.PutExportFile(zipFilePath)
	if err != nil {
		reason := "上传阿里云oss失败：" + err.Error()
		expLog.SaveFailReason(reason)
		return
	}
	os.RemoveAll(zipFilePath)

	// 6. 修改任务状态，写文件
	expLog.Status = mdb.ExportLog_status_succ
	dao.MDB.Model(&expLog).Select("status").Updates(expLog)
	// 创建文件数据
	expFile := &mdb.ExportFile{
		HashKey: expLog.HashKey,
		Path:    objname,
		Type:    expLog.ExtType,
	}
	res := dao.MDB.Create(expFile)
	if res.Error != nil {
		glog.Error("exportfile insert err", "", res.Error.Error())
		return
	}
	// w.notify(expLog.Callback, taskID)
	new(cal.SourceHTTP).Notify(expLog.Callback, taskID)
	log.Print(taskID, "任务完成")
}
