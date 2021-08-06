package cal

import (
	"errors"
	"export-server/models/dao/mdb"
	"export-server/pkg/glog"
	"fmt"
	"net/url"
	"strings"
	"sync"

	request "gitee.com/smallcatx0/gequest"
	"github.com/tidwall/gjson"
)

type SourceHTTP struct {
	Cli *request.Core
}

func NewSHttp() *SourceHTTP {
	shttp := &SourceHTTP{}
	shttp.Cli = request.New("export-servers", "", 10000)
	return shttp
}

func NewSHttpWCli(cli *request.Core) *SourceHTTP {
	return &SourceHTTP{
		Cli: cli,
	}
}

func (r *SourceHTTP) BuildHttpSource(param HttpParam) request.Core {
	// 拷贝param 并发安全
	copyParam := make(map[string]interface{}, len(param.Param))
	for k, v := range param.Param {
		copyParam[k] = v
	}

	copyParam["page"] = param.Page
	if _, ok := copyParam["limit"]; !ok {
		copyParam["limit"] = 50
	}
	method := strings.ToLower(param.Method)
	req := request.New("export-servers", "", 30000).
		SetMethod(method).
		SetUri(param.Url).
		AddHeaders(param.Header)
	switch method {
	case "post":
		req.SetJson(copyParam)
	case "get":
		q := url.Values{}
		for k, v := range copyParam {
			q.Add(k, fmt.Sprintf("%v", v))
		}
		req.SetQuery(q)
	}
	return *req
}

func (r *SourceHTTP) PaseRes(res *request.Response) (totalpage int, lists string, err error) {
	bodyStr, err := res.ToString()
	bodyJson := gjson.Parse(bodyStr)
	totalpage = int(bodyJson.Get("data.meta.pagination.total_pages").Int())
	lists = bodyJson.Get("data.data").String()
	return
}

func (r *SourceHTTP) GetHttpSource(
	param HttpParam,
) (totalpage int, page int, lists string, err error) {
	req := r.BuildHttpSource(param)
	res, err := req.SendRtry(3)
	if err != nil {
		glog.ErrorT("http request err", "", err, req.String())
		return
	}
	bodyStr, err := res.ToString()
	bodyJson := gjson.Parse(bodyStr)
	statusCode := int(bodyJson.Get("status_code").Int())
	if statusCode != 0 {
		err = errors.New(bodyJson.String())
		return
	}
	totalpage = int(bodyJson.Get("data.pagetag.total_page").Int())
	page = int(bodyJson.Get("data.pagetag.page").Int())
	lists = bodyJson.Get("data.data").String()
	return
}

type HttpParam struct {
	Page        int
	Url, Method string
	Header      map[string]string
	Param       map[string]interface{}
}

var wg sync.WaitGroup

type Resp struct {
	Page  int
	Lists string
	Err   error
}

func (r *SourceHTTP) MultPageReq(param ...HttpParam) (respCh chan Resp) {
	wg.Add(len(param))
	respCh = make(chan Resp, len(param)+1)
	for _, aparam := range param {
		go func(param HttpParam, ch chan Resp) {
			defer wg.Done()
			_, page, lists, err := r.GetHttpSource(param)
			ch <- Resp{
				Page:  page,
				Lists: lists,
				Err:   err,
			}
		}(aparam, respCh)
	}
	wg.Wait()
	close(respCh)
	return
}

func (r *SourceHTTP) Notify(url, taskID string) {
	if url == "" {
		// 回调地址为空 直接跳过
		return
	}
	// 查询结果
	taskdetail, _ := new(mdb.ExportLog).Detail(taskID)
	r.Cli.Clear().SetUri(url).
		SetMethod("post").
		SetJson(taskdetail).
		SendRtry(5)
}
