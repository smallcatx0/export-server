package routes

import (
	v1 "export-server/controller/v1"

	"github.com/gin-gonic/gin"
)

func registeRoute(router *gin.Engine) {
	router.GET("/demo/page", v1.PageDemo)

	router.POST("/v1/export/http", v1.ExportSHttp)
	router.POST("/v1/export/raw", v1.ExportSRaw)
	router.GET("/v1/export/detail", v1.ExportDetail)
	router.GET("/v1/export/history", v1.ExportHistory)
}
