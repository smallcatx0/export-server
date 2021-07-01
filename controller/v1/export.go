package v1

import (
	"export-server/middleware/httpmd"
	"export-server/models/page"
	"export-server/pkg/exception"
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
	list := make([]interface{}, 1, 10)
	list[0] = map[string]string{"ID": "编号", "Name": "姓名", "Age": "年龄"}
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

// ExportSHttp http接口导出excel
func ExportSHttp(c *gin.Context) {
	param := valid.ExpSHttpParam{}
	err := valid.BindAndCheck(c, &param)
	if err != nil {
		r.Fail(c, err)
		return
	}
	exportServ := new(page.ExportServ)
	data, err := exportServ.HandelSHttp(c, &param)
	if err != nil {
		r.Fail(c, err)
		return
	}
	r.Succ(c, data)
}

// ExportSRaw 源数据导出excel
func ExportSRaw(c *gin.Context) {
	param := valid.ExpSRawParam{}
	err := valid.BindAndCheck(c, &param)
	if err != nil {
		r.Fail(c, err)
		return
	}
	exportServ := new(page.ExportServ)
	data, err := exportServ.HandelSRaw(c, &param)
	if err != nil {
		r.Fail(c, err)
		return
	}
	r.Succ(c, data)

}

func ExportDetail(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		r.Fail(c, exception.ParamInValid("key 不能为空"))
		return
	}

	data, err := new(page.ExportServ).Detail(key)
	if err != nil {
		r.Fail(c, err)
		return
	}
	r.Succ(c, data)

}
