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

func TestXY2Pos(t *testing.T) {
	assert := assert.New(t)
	cases := map[string][2]int{
		"A1":   [2]int{1, 1},
		"B2":   [2]int{2, 2},
		"I2":   [2]int{9, 2},
		"AF98": [2]int{32, 98},
		"AI3":  [2]int{35, 3},
		"AO1":  [2]int{41, 1},
		"BI70": [2]int{61, 70},
		"BU7":  [2]int{73, 7},
		"AN56": [2]int{40, 56},
	}
	for except, param := range cases {
		p := excel.Pos{X: param[0], Y: param[1]}
		act := p.String()
		assert.Equal(except, act, param)
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
	assert := assert.New(t)
	var listJson = `[{"ID":0,"Name":"Demo","Age":15,"Time":1231231},{"ID":0,"Name":"Demo","Age":15,"Time":1231231},{"ID":0,"Name":"Demo","Age":15,"Time":1231231}]`
	jsonOj := gjson.Parse(listJson).Array()
	keys := excel.JsonKeys(jsonOj[0].String())
	listMap := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(listJson), &listMap)
	marr := excel.List2Arrs(listMap, keys)
	assert.EqualValues(marr[0][0], 0)
	assert.EqualValues(marr[1][1], "Demo")
	assert.EqualValues(marr[2][2], 15)
}

func TestJson2Excel(t *testing.T) {
	input, _ := ioutil.ReadFile("./.assets/list-500.json")
	excelOj := excel.NewExcelRecorder("/tmp/outExcel/list-500x5.xlsx")
	p := excelOj.JsonListWrite(excel.Pos{X: 1, Y: 1}, string(input), true)
	for i := 0; i < 5; i++ {
		p = excelOj.JsonListWrite(p, string(input), false)
		fmt.Println(p)
	}
	excelOj.Save()
}

// 5w行分多文件 5000行一个文件
func TestWritPaging(t *testing.T) {
	pageLimit := 5000
	input, _ := ioutil.ReadFile("./.assets/list-5000.json")
	oj := excel.NewExcelRecorderPage("/tmp/outExcel/list-5w-%d.xlsx", pageLimit)
	p := excel.Pos{X: 1, Y: 1}
	for i := 0; i < 10; i++ {
		p = oj.WritePagpenate(p, string(input), "")
	}
}

func BenchmarkJson2Excel(b *testing.B) {
	input, _ := ioutil.ReadFile("./.assets/list-5.json")
	lines := gjson.Parse(string(input)).Array()
	htable := lines[0].String()
	excelOj := excel.NewExcelRecorder("/tmp/outExcel/BenchmarkJson2Excel.xlsx")
	p := excel.Pos{X: 1, Y: 5000}
	keys := excel.JsonKeys(htable)
	excelOj.JsonWrite(p, keys, lines[0].String())

	for i := 0; i < b.N; i++ {
		p = excelOj.JsonWrite(p, keys, lines[1].String())
	}
	excelOj.Save()
}
