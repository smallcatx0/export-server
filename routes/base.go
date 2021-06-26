package routes

import (
	C "export-server/controller"
	"export-server/pkg/conf"

	"github.com/gin-gonic/gin"
)

// Register http路由总入口
func Register(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		v := conf.AppConf.GetString("base.describe")
		c.String(200, v)
	}) // version
	r.GET("/healthz", C.Healthz)
	r.GET("/ready", C.Ready)
	r.GET("/reload", C.ReloadConf)
	r.GET("/test", C.Test)
	registeRoute(r)
}
