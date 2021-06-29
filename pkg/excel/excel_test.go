package excel_test

import (
	"encoding/json"
	"export-server/pkg/excel"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestX2column(t *testing.T) {
	assert := assert.New(t)
	columnMap := map[int]string{
		1: "A", 2: "B", 9: "I", 32: "AF", 35: "AI", 41: "AO", 61: "BI", 73: "BU",
	}
	for k, act := range columnMap {
		v := excel.X2col(k)
		assert.Equal(v, act, k)
	}

}

func TestJsonKeys(t *testing.T) {
	assert := assert.New(t)
	var input = `{"ID":0,"Name":"Demo","Age":15,"Time":1231231}`
	expect := []string{"ID", "Name", "Age", "Time"}
	for i := 0; i < 1000; i++ {
		tmp := excel.JsonKeys(input)
		assert.Equal(tmp, expect)
	}
}

func TestJson2Arr(t *testing.T) {
	// assert := assert.New(t)
	var listJson = `[{"ID":0,"Name":"Demo","Age":15,"Time":1231231},{"ID":0,"Name":"Demo","Age":15,"Time":1231231},{"ID":0,"Name":"Demo","Age":15,"Time":1231231}]`
	jsonOj := gjson.Parse(listJson).Array()
	keys := excel.JsonKeys(jsonOj[0].String())
	listMap := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(listJson), &listMap)
	marr := excel.List2Arrs(listMap, keys)
	fmt.Println(marr)
}

func TestExcelWritLines(t *testing.T) {
	exOj := excel.NewExcelRecorder("./out.xlsx")
	listData := [][]interface{}{
		[]interface{}{"ID", "Name", "Age"},
		[]interface{}{1, "kkkk", 30},
		[]interface{}{1, "tttt", 11},
	}
	exOj.WriteLines(excel.Pos{X: 1, Y: 1}, listData)
	exOj.Save()
}

func TestJson2Excel(t *testing.T) {
	input, _ := ioutil.ReadFile("./.assets/list-5000.json")
	excelOj := excel.NewExcelRecorder("./.assets/out.xlsx")
	p, _ := excelOj.JsonListWrite(excel.Pos{X: 1, Y: 1}, string(input), true)
	for i := 0; i < 10; i++ {
		p, _ = excelOj.JsonListWrite(p, string(input), false)
		fmt.Println(p)
	}
	excelOj.Save()
}

func TestJson2Excel2(t *testing.T) {
	input, _ := ioutil.ReadFile("./.assets/list-5000.json")
	oj := excel.NewExcelRecorder("./.assets/out.xlsx")
	oj.JsonListWrite2(excel.Pos{X: 1, Y: 1}, string(input), true)
	oj.Save()
}

func BenchmarkJson2Excel(b *testing.B) {
	input, _ := ioutil.ReadFile("./.assets/list-500.json")
	excelOj := excel.NewExcelRecorder("./.assets/out.xlsx")
	for i := 0; i < b.N; i++ {
		excelOj.JsonListWrite(excel.Pos{X: 1, Y: 1}, string(input), true)
	}
	excelOj.Save()
}

func BenchmarkJson2Excel2(b *testing.B) {
	input, _ := ioutil.ReadFile("./.assets/list-500.json")
	excelOj := excel.NewExcelRecorder("./.assets/out.xlsx")
	for i := 0; i < b.N; i++ {
		excelOj.JsonListWrite2(excel.Pos{X: 1, Y: 1}, string(input), true)
	}
	excelOj.Save()
}
