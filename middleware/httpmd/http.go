package httpmd

import (
	"encoding/json"
	"export-server/pkg/conf"
	glog "export-server/pkg/glog"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// RequestIDKey 唯一请求id
const RequestIDKey = "x-b3-traceid"

// 日志记录白名单
var LogWrite = map[string]bool{
	"/healthz": true,
	"/ready":   true,
}

// SetHeader 设置header
func SetHeader(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	requestID := c.GetHeader(RequestIDKey)
	if requestID == "" {
		c.Set(RequestIDKey, uuid.NewV4().String())
	} else {
		c.Set(RequestIDKey, requestID)
	}
}

// ReqLog 记录全量请求日志
func ReqLog(c *gin.Context) {
	if conf.IsDebug() {
		path := c.Request.RequestURI
		requestData, _ := c.Copy().GetRawData()
		header, _ := json.Marshal(c.Request.Header)
		if _, ok := LogWrite[path]; ok {
			// 白名单不记录
			return
		}
		glog.Debug("request_log",
			c.GetString(RequestIDKey),
			path,
			string(requestData),
			string(header),
		)
	}
}
