package v1

import (
	"export-server/middleware/httpmd"
	"export-server/valid"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// PageDemo 分页接口demo
func PageDemo(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	pageinfo := httpmd.Pagination{}
	pageinfo.Format(page, limit, 500)
	list := make([]interface{}, 0, 10)

	for i := 0; i < pageinfo.Limit; i++ {
		tmp := struct {
			ID   int
			Name string
			Age  int
		}{pageinfo.Offset + i, "Demo", 15}
		list = append(list, tmp)
	}
	data := map[string]interface{}{
		"pagetag": pageinfo,
		"data":    list,
	}
	r.Succ(c, data)
}

func Export(c *gin.Context) {
	param := &valid.ExportParam{}
	err := valid.BindAndCheck(c, param)
	if err != nil {
		r.Fail(c, err)
		return
	}
	switch strings.ToLower(param.SourceType) {
	case "http":
		// 1. 获取参数哈希存日志
		// 2. 准备参数丢任务队列中
		// 3. 返回参数哈希
	}

}
