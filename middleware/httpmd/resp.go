package httpmd

import (
	"net/http"
	"strings"

	"export-server/pkg/conf"
	"export-server/pkg/exception"
	glog "export-server/pkg/glog"
	"export-server/pkg/helper"

	"github.com/gin-gonic/gin"
)

// Resp 封装响应体
type Resp struct{}

// 响应体
type responseData struct {
	ErrCode   uint32      `json:"errcode"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id"`
}

// 分页规范
type Pagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
	Offset    int `json:"-"`
}

// Format 格式化
func (p *Pagination) Format(page, limit, total int) {
	p.Total = total
	p.Limit = helper.GetDefInt(limit, 10)
	p.TotalPage = p.Total / p.Limit
	if (p.Total / p.Limit) == 0 {
		p.TotalPage = 1
	}
	switch {
	case page < 1:
		p.Page = 1
	case page >= 1 && page <= p.TotalPage:
		p.Page = page
	case page > p.TotalPage:
		p.Page = p.TotalPage
	}
	p.Offset = (p.Page - 1) * p.Limit
}

// Succ 成功返回
func (r *Resp) Succ(c *gin.Context, data interface{}, msg ...string) {
	rr := new(responseData)
	rr.ErrCode = 0
	if len(msg) == 0 {
		rr.Msg = exception.ErrNos[rr.ErrCode]
	} else {
		rr.Msg = strings.Join(msg, ",")
	}
	rr.Data = data
	rr.RequestID = c.GetString(RequestIDKey)
	c.JSON(http.StatusOK, &rr)
}

// Fail 失败返回
func (r *Resp) Fail(c *gin.Context, err error) {
	var httpState int
	rr := new(responseData)
	switch e := err.(type) {
	case *exception.Exception:
		rr.ErrCode = e.Code
		rr.Msg = e.Msg
		httpState = e.HTTPCode
	default:
		// 记录日志
		if conf.Env() == "dev" {
			rr.Msg = err.Error()
		} else {
			glog.Error(c.Request.RequestURI, err.Error())
			rr.Msg = "服务错误"
		}
		rr.ErrCode = 400
		httpState = 400
	}
	rr.RequestID = c.GetString(RequestIDKey)
	c.JSON(httpState, &rr)
}

func (r *Resp) SuccJsonRaw(c *gin.Context, data string) {
	format := `{"errcode":%d,"msg":"%s","data":%s,"request_id":"%s"}`
	requestId := c.GetString(RequestIDKey)
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(200, format, 200, "操作成功", data, requestId)
}
