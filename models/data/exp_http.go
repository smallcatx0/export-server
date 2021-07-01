package data

import (
	"encoding/json"
	"export-server/bootstrap/global"
	"export-server/models/dao"
	"export-server/models/dao/mdb"
	"export-server/models/dao/rdb"
	"export-server/pkg/conf"
	"export-server/pkg/excel"
	"export-server/pkg/glog"
	"export-server/pkg/helper"
	"export-server/valid"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	request "gitee.com/smallcatx0/gequest"
	"github.com/tidwall/gjson"
)

type HttpWorker struct {
	Tasks  *rdb.Mq
	Cli    *request.Core
	taskCh chan *rdb.ExportTask
}

func (w *HttpWorker) Run(pool int) {
	// 单消费端 多任务执行
	w.Tasks = &rdb.Mq{Key: global.TaskHttpKey}
	w.Cli = request.New("export-server", "", 3000).Debug(conf.IsDebug())
	// 缓冲区越大，程序宕机后丢消息越多
	w.taskCh = make(chan *rdb.ExportTask, 20)
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
	atask := <-w.taskCh
	// 1. 数据库中查询任务详情
	expLog := mdb.ExportLog{}
	result := dao.MDB.Where("hash_key=?", atask.TaskID).First(&expLog)
	if result.Error != nil {
		glog.Error("TaskNotFund hash_key=" + atask.TaskID)
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
		err := expLog.SaveFailReason(reason)
		if err != nil {
			glog.Error("udate export_log err", "", err.Error())
		}
		return
	}
	// TODO: 并行请求 让单个任务更快完成
	// 但也要保证顺序

	// 2. 获取数据源的数据 -> 3. 写入excel
	totalPage, list := w.getSource(requestParam, 1) // 第一页
	excelTmpPath := conf.AppConf.GetString("storage.outexcel_tmp")
	filename := expLog.Title + "-%d." + expLog.ExtType
	excelw := excel.NewExcelRecorderPage(path.Join(excelTmpPath, atask.TaskID, filename), 200)
	p := excelw.WritePagpenate(excel.Pos{X: 1, Y: 1}, list, "", true)
	for i := 2; i <= totalPage; i++ { // 循环获取剩下页
		log.Printf("开始抓取%d页\n", i)
		_, list = w.getSource(requestParam, i)
		p = excelw.WritePagpenate(p, list, "", false)
	}
	excelw.Save()

	// 4. 压缩文件夹 并删除源文件
	zipFilePath := path.Join(excelTmpPath, atask.TaskID+".zip")
	taskDir := path.Join(excelTmpPath, atask.TaskID)
	helper.FolderZip(taskDir, zipFilePath)
	os.RemoveAll(taskDir)

	// 5. 上传云 OOS -> 删除本地文件
	// ...

	// 6. 修改任务状态，写文件
	expLog.Status = mdb.ExportLog_status_succ
	dao.MDB.Model(&expLog).Select("status").Updates(expLog)
	// 创建文件数据
	expFile := &mdb.ExportFile{
		HashKey: expLog.HashKey,
		Path:    zipFilePath,
		Type:    expLog.ExtType,
	}
	res := dao.MDB.Create(expFile)
	if res.Error != nil {
		glog.Error("exportfile insert err", "", res.Error.Error())
		return
	}
	// TODO: 请求回调通知

	log.Print("任务完成 ", atask.TaskID)
}

func (w *HttpWorker) getSource(reqParam valid.SourceHTTP, page int) (totalPage int, lists string) {
	// 分页逐个请求
	reqParam.Param["page"] = page
	method := strings.ToLower(reqParam.Method)
	req := w.Cli.SetMethod(method).
		SetUri(reqParam.URL).
		AddHeaders(reqParam.Header)
	switch method {
	case "post":
		req.SetJson(reqParam.Param)
	case "get":
		q := url.Values{}
		for k, v := range reqParam.Param {
			q.Add(k, fmt.Sprintf("%v", v))
		}
		req.SetQuery(q)
	}
	res, err := req.Send()
	if err != nil {
		// TODO: httpq请求失败需要重试
		glog.ErrorT("http request err", "", err, reqParam)
		return
	}
	bodyStr, err := res.ToString()
	if err != nil {
		glog.Error("http respons body read err", "", err.Error())
		return
	}
	bodyJson := gjson.Parse(bodyStr)
	totalPage = int(bodyJson.Get("data.pagetag.total_page").Int())
	lists = bodyJson.Get("data.data").String()
	return
}
