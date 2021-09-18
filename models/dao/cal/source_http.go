package cal

import (
	"encoding/json"
	"errors"
	"export-server/models/dao/mdb"
	"export-server/pkg/conf"
	"export-server/pkg/glog"
	"fmt"
	"net/url"
	"strings"

	request "github.com/smallcatx0/gequest"
	"github.com/tidwall/gjson"
)

type HttpParam struct {
	Page        int
	Url, Method string
	Header      map[string]string
	Param       map[string]interface{}
}

type Logger struct{}

func (Logger) Print(v ...interface{}) {
	msgs := make([]string, 1, len(v))
	for _, a := range v {
		switch one := a.(type) {
		case fmt.Stringer:
			msgs = append(msgs, one.String())
		case string:
			msgs = append(msgs, one)
		default:
			jsonv, _ := json.Marshal(a)
			msgs = append(msgs, string(jsonv))
		}
	}
	glog.Debug("request_log", msgs...)
}

type SourceHTTP struct {
	Cli *request.Core
}

func NewSourceHTTP() *SourceHTTP {
	s := &SourceHTTP{}
	s.Cli = request.New("export-server", "", 1000)
	s.Cli.Debug(conf.IsDebug())
	s.Cli.SetLoger(Logger{})
	return s
}

// 回调通知结果
func (s *SourceHTTP) Notify(url, taskID string) {
	if url == "" {
		// 回调地址为空 直接跳过
		return
	}
	// 查询结果
	taskdetail, _ := new(mdb.ExportLog).Detail(taskID)
	s.Cli.Clear().SetUri(url).
		SetMethod("post").
		SetJson(taskdetail).
		SendRtry(5)
}

// 构建请求
func (s *SourceHTTP) BuildReq(param *HttpParam) request.Core {
	// 拷贝map 使其并发安全
	query := make(map[string]interface{}, len(param.Param))
	for k, v := range param.Param {
		query[k] = v
	}
	query["page"] = param.Page
	if _, ok := query["limit"]; !ok {
		query["limit"] = 50
	}
	method := strings.ToUpper(param.Method)
	// TODO: 这里应该新建一个 request.Core 对象
	req := s.Cli.SetMethod(method).
		SetUri(param.Url).
		SetHeaders(param.Header)
	switch method {
	case "POST":
		req.SetJson(query)
	case "GET":
		q := url.Values{}
		for k, v := range query {
			q.Add(k, fmt.Sprintf("%v", v))
		}
		req.SetQuery(q)
	}
	return *req
}

// 解析结果
func (s *SourceHTTP) PaseResponse(r *request.Response) (
	page int,
	totalPage int,
	lists string,
	err error,
) {
	bodyStr, err := r.ReadAll()
	bodyParse := gjson.ParseBytes(bodyStr)
	errCode := int(bodyParse.Get("errcode").Int())
	if errCode != 0 {
		err = errors.New(bodyParse.String())
		return
	}
	totalPage = int(bodyParse.Get("data.pagetag.total_page").Int())
	page = int(bodyParse.Get("data.pagetag.page").Int())
	lists = bodyParse.Get("data.data").String()
	return
}

// 请求第一页
func (s *SourceHTTP) FirstPage(param *HttpParam) (
	page int,
	totalPage int,
	lists string,
	err error,
) {
	req := s.BuildReq(param)
	r, err := req.SendRtry(5)
	if err != nil {
		glog.Error("send http request err or tan 5 times", "", err.Error(), req.String())
		return
	}
	return s.PaseResponse(r)
}

// 并发请求
func (s *SourceHTTP) BatchRequest(params ...HttpParam) (
	lists []string,
	err error,
) {
	reqs := make([]*request.Core, 0, len(params))
	lists = make([]string, 0, len(params))
	// 构建请求
	for _, param := range params {
		aReq := s.BuildReq(&param)
		reqs = append(reqs, &aReq)
	}
	// 并发请求
	multRes := request.MultRequest(5, reqs...)
	for _, one := range multRes {
		var alist string
		if err = one.Err; err != nil {
			// 有任何一个错误即返回
			glog.Error("send http request err", "", err.Error(), one.Core.String())
			return
		}
		// 解析响应
		_, _, alist, err = s.PaseResponse(one.Core.Response())
		if err != nil {
			glog.Error("pase http response err", "", err.Error(), one.Core.String())
			return
		}
		lists = append(lists, alist)
	}
	return
}
