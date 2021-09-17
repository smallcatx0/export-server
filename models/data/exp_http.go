package data

import (
	"encoding/json"
	"export-server/bootstrap/global"
	"export-server/models/dao"
	cal "export-server/models/dao/cal"
	"export-server/models/dao/mdb"
	"export-server/models/dao/rdb"
	"export-server/pkg/conf"
	"export-server/pkg/excel"
	"export-server/pkg/glog"
	"export-server/pkg/helper"
	"export-server/valid"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/golang-module/carbon"
)

type HttpWorker struct {
	Tasks  *rdb.Mq
	taskCh chan *rdb.ExportTask
	req    *cal.SourceHTTP
}

func (w *HttpWorker) Run(pool int) {
	// 单消费端 多任务执行
	w.Tasks = &rdb.Mq{Key: global.TaskHttpKey}
	// 缓冲区越大，程序宕机后丢消息越多
	w.taskCh = make(chan *rdb.ExportTask, 20)
	w.req = cal.NewSourceHTTP()
	log.Print("httpWorker pool=", pool)
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
func (w *HttpWorker) startWorker(pool int) {
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

// work 处理单个任务
func (w *HttpWorker) work() {
	currTask := <-w.taskCh
	taskID := currTask.TaskID
	st := carbon.Now()
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
	requestParam := valid.SourceHTTP{}
	err := json.Unmarshal([]byte(expLog.Param), &requestParam)
	if err != nil {
		reason := "参数解析失败：" + err.Error()
		expLog.SaveFailReason(reason)
		w.req.Notify(expLog.Callback, taskID)
		return
	}

	// 2. 获取数据源的数据 -> 3. 写入excel
	baseParam := &cal.HttpParam{
		Page:   1,
		Url:    requestParam.URL,
		Method: requestParam.Method,
		Header: requestParam.Header,
		Param:  requestParam.Param,
	}
	// 带上此次任务ID 方便日志追踪
	if baseParam.Header == nil {
		baseParam.Header = make(map[string]string)
	}
	baseParam.Header["xt-export-taskId"] = taskID
	// 第一次请求 获取分页信息
	page, totalPage, lists, err := w.req.FirstPage(baseParam)
	log.Printf("[%s] 抓取到(%d/%d)页\n", taskID, page, totalPage)
	glog.InfoF("抓取到 (%d/%d)页", taskID, page, totalPage)
	if err != nil {
		reason := "获取数据源失败：" + err.Error()
		expLog.SaveFailReason(reason)
		w.req.Notify(expLog.Callback, taskID)
		return
	}
	// 读配置
	excelTmpPath := conf.AppConf.GetString("tmp_storage.outexcel_tmp") // excel 临时文件目录
	maxlines := conf.AppConf.GetInt("excel_maxlines") + 1              // excel 最大行数

	conn := 1 // 默认无并发
	if requestParam.Conn != 0 {
		conn = requestParam.Conn
	}
	filename := expLog.Title + "-%d." + expLog.ExtType
	excelw := excel.NewExcelRecorderPage(
		path.Join(excelTmpPath, taskID, filename),
		maxlines,
	)
	// 写入第一页
	p := excelw.WritePagpenate(excel.Pos{X: 1, Y: 1}, lists, "", true)
	page += 1
	var end bool
	for { // 获取剩下页
		params := make([]cal.HttpParam, 0, conn)
		for j := 0; j < conn; j++ {
			if page > totalPage {
				end = true
				break
			}
			baseParam.Page = page
			params = append(params, *baseParam)
			log.Printf("[%s] 开始抓取(%d/%d)页\n", taskID, page, totalPage)
			glog.InfoF("开始抓取(%d/%d)页", taskID, page, totalPage)
			page += 1
		}
		lists, err := w.req.BatchRequest(params...)
		if err != nil {
			reason := "获取数据源失败：" + err.Error()
			expLog.SaveFailReason(reason)
			w.req.Notify(expLog.Callback, taskID)
		}
		// 结果写入
		for _, alist := range lists {
			p = excelw.WritePagpenate(p, alist, "", false)
		}
		if end {
			break
		}
	}
	excelw.Save()

	// 4. 压缩文件夹 并删除源文件
	zipFilePath := path.Join(excelTmpPath, taskID+".zip")
	taskDir := path.Join(excelTmpPath, taskID)
	helper.FolderZip(taskDir, zipFilePath)
	os.RemoveAll(taskDir)

	// 5. 文件持久化（阿里云oss、本地）-> 删除本地文件
	objname, err := dao.FS.Put(zipFilePath)
	if err != nil {
		reason := "文件持久化失败" + err.Error()
		expLog.SaveFailReason(reason)
		return
	}
	os.RemoveAll(zipFilePath)

	// 6. 修改任务状态
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
	w.req.Notify(expLog.Callback, taskID)
	dt := carbon.Now().DiffInSecondsWithAbs(st)
	log.Printf("[%s] 任务完成 耗时%ds", taskID, dt)
	glog.InfoF("任务完成 耗时%ds", taskID, dt)
}
