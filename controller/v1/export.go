package v1

import (
	"export-server/middleware/httpmd"
	"export-server/models/page"
	"export-server/valid"
	"strconv"

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

// 数据导出
func Export(c *gin.Context) {
	param := &valid.ExportParam{}
	err := valid.BindAndCheck(c, param)
	if err != nil {
		r.Fail(c, err)
		return
	}
	exportServ := new(page.ExportServ)
	data, err := exportServ.Handel(c, param)
	if err != nil {
		r.Fail(c, err)
		return
	}
	r.Succ(c, data)
}
