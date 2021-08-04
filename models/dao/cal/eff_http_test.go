package cal_test

import (
	cal "export-server/models/dao/Cal"
	"log"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildHttp(t *testing.T) {
	request := cal.NewSHttp()
	param := cal.HttpParam{
		Page: 1,
		Url:  "http://127.0.0.1:8080/demo/page", Method: "get",
		Header: map[string]string{},
		Param:  map[string]interface{}{"per_page": 5},
	}
	req1 := request.BuildHttpSource(param)
	param.Page += 1
	req2 := request.BuildHttpSource(param)
	param.Page += 1
	req3 := request.BuildHttpSource(param)

	res, _ := req1.SendRtry(3)
	_, lists1, _ := request.PaseRes(res)
	log.Println(lists1)
	res, _ = req2.SendRtry(3)
	_, lists2, _ := request.PaseRes(res)
	log.Println(lists2)
	res, _ = req3.SendRtry(3)
	_, lists3, _ := request.PaseRes(res)
	log.Println(lists3)
}

func TestMultPageReq(t *testing.T) {
	request := cal.NewSHttp()
	param := cal.HttpParam{
		Page: 1,
		Url:  "http://127.0.0.1:8080/demo/page", Method: "get",
		Header: map[string]string{},
		Param:  map[string]interface{}{"current_page": 1, "per_page": 5},
	}
	params := make([]cal.HttpParam, 0)
	params = append(params, param)
	param.Page += 1
	params = append(params, param)
	param.Page += 1
	params = append(params, param)
	param.Page += 1
	params = append(params, param)
	param.Page += 1
	params = append(params, param)
	param.Page += 1
	params = append(params, param)
	respCh := request.MultPageReq(params...)

	for one := range respCh {
		log.Print(one)
	}
	log.Print("end")
}

func TestGetHttpSource(t *testing.T) {
	ass := assert.New(t)
	bin := 3
	request := cal.NewSHttp()
	param := cal.HttpParam{
		Page: 1,
		Url:  "http://127.0.0.1:8080/demo/page", Method: "get",
		Header: map[string]string{},
		Param:  map[string]interface{}{"per_page": 5},
	}
	// 先单独请求一次
	tpage, _, lists, err := request.GetHttpSource(param)
	ass.NoError(err)
	log.Print(lists)
	i := 2
	var end bool
	for {
		params := make([]cal.HttpParam, 0)
		// 并发分批请求
		for j := 0; j < bin; j++ {
			if i > tpage {
				end = true
				break
			}
			param.Page = i
			params = append(params, param)
			i += 1
		}
		respCh := request.MultPageReq(params...)
		// 将数据处理为有序
		res := make(map[int]string, bin)
		keys := make([]int, 0, bin)
		for one := range respCh {
			res[one.Page] = one.Lists
			keys = append(keys, one.Page)
		}
		sort.Ints(keys)
		log.Print(keys)
		log.Print("------------------------------")
		if end {
			break
		}
	}

}
