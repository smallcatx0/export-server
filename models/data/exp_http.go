package data

import (
	"encoding/json"
	"export-server/bootstrap/global"
	"export-server/models/dao"
	"export-server/models/dao/mdb"
	"export-server/models/dao/rdb"
	"export-server/pkg/glog"
	"export-server/valid"
	"fmt"
	"log"
	"net/url"
	"strings"

	request "gitee.com/smallcatx0/gequest"
	"github.com/tidwall/gjson"
)

type HttpWorker struct {
	Tasks *rdb.Mq
	Cli   *request.Core
	// Doubt: 携程间消息通道是直接沿用消息队列中的结构 还是只留hash_key
	// 沿用消息队列的结构感觉稍微有点耦合、但是自留hash_key 有损失了可扩展性
	taskCh chan *rdb.ExportTask
}

func (w *HttpWorker) Run(pool int) {
	// 单消费端 多任务执行
	w.Tasks = &rdb.Mq{Key: global.TaskHttpKey}
	w.Cli = request.New("export-server", "", 3000).Debug(true)
	// 缓冲区越大，程序宕机后丢消息越多
	w.taskCh = make(chan *rdb.ExportTask, 5)
	log.Print(pool)
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

// startWorker 启动工作携程
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
	// 数据库中查询任务详情
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
		glog.Error("json.Unmarshal err " + err.Error())
	}
	// TODO: 并行请求 让单个任务更快完成
	// 获取数据源的数据
	totalPage, list := w.GetSource(requestParam, 1) // 第一页
	// 写入excl

	// for i := 2; i < totalPage; i++ {
	// 	_, currlist := w.GetSource(requestParam, i)

	// }
	// 写入excl文件
	log.Print(list, totalPage)
}

func (w *HttpWorker) GetSource(reqParam valid.SourceHTTP, page int) (totalPage int, list []map[string]interface{}) {
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
		glog.Error("http request err", "", err.Error())
		return
	}
	bodyStr, err := res.ToString()
	if err != nil {
		glog.Error("http respons body read err", "", err.Error())
		return
	}
	bodyJson := gjson.Parse(bodyStr)
	totalPage = int(bodyJson.Get("data.pagetag.total_page").Int())
	// 循环获取数据，
	listRaw := bodyJson.Get("data.data").String()
	list = make([]map[string]interface{}, 0, 10)
	err = json.Unmarshal([]byte(listRaw), &list)
	if err != nil {
		glog.Error("http数据源json 解析失败", "", listRaw)
		return
	}
	return
}
