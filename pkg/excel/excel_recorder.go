package excel

import (
	"export-server/pkg/helper"
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tidwall/gjson"
)

type ExcelRecorder struct {
	FilePath         string
	ExcelFp          *excelize.File
	Sheet            string
	Limit            int
	Page             int
	FileNameTemplate string
}

func NewExcelRecorder(path string) *ExcelRecorder {
	instance := &ExcelRecorder{
		FilePath: path,
		ExcelFp:  excelize.NewFile(),
		Sheet:    "Sheet1",
	}
	return instance
}

func NewExcelRecorderPage(template string, limit int) *ExcelRecorder {
	ins := &ExcelRecorder{
		ExcelFp:          excelize.NewFile(),
		Sheet:            "Sheet1",
		Page:             1,
		Limit:            limit,
		FileNameTemplate: template,
	}
	return ins
}

// JsonListWrite 将json列表写入excel中
func (e *ExcelRecorder) JsonListWrite(p Pos, jsonStr string, isFirst bool) Pos {
	result := gjson.Parse(jsonStr)
	keys := make([]string, 0)
	for i, line := range result.Array() {
		if i == 0 {
			keys = JsonKeys(line.String())
			if !isFirst {
				// 不写首行数据
				continue
			}
		}
		p = e.JsonWrite(p, keys, line.String())
	}
	return p
}

// JsonWrite 单条json 按照keys的键顺序写入excel
// ajson格式 `{"key1":"val1", "key2":"val2"}`
func (e *ExcelRecorder) JsonWrite(p Pos, keys []string, ajson ...string) Pos {
	for _, line := range ajson {
		aline := gjson.Parse(line).Value().(map[string]interface{})
		lineValues := helper.Map2Arr(aline, keys)
		e.ExcelFp.SetSheetRow(e.Sheet, p.String(), &lineValues)
		p.Y += 1
	}
	return p
}

// Save 保存
func (e *ExcelRecorder) Save() error {
	if e.FileNameTemplate != "" {
		e.FilePath = fmt.Sprintf(e.FileNameTemplate, e.Page)
	}
	helper.TouchDir(e.FilePath)
	if err := e.ExcelFp.SaveAs(e.FilePath); err != nil {
		return err
	}
	return nil
}

func List2Arrs(lines []map[string]interface{}, keys []string) [][]interface{} {
	ret := make([][]interface{}, 0, len(lines))
	for _, line := range lines {
		tmp := helper.Map2Arr(line, keys)
		ret = append(ret, tmp)
	}
	return ret
}

// 分页写入
func (e *ExcelRecorder) WritePagpenate(p Pos, lineJson string, htable string, isFirst bool) Pos {
	lines := gjson.Parse(lineJson).Array()
	if htable == "" {
		htable = lines[0].String()
		lines = lines[1:]
	}
	keys := JsonKeys(htable)
	// 第一次调用 需要写表头
	if isFirst {
		p = e.JsonWrite(p, keys, htable)
	}
	for _, aline := range lines {
		if p.Y >= e.Limit {
			// 保存文件
			e.Save()
			// 新生成文件
			e.ExcelFp = excelize.NewFile()
			e.Page += 1
			p.Y = 1
			// 写表头
			p = e.JsonWrite(p, keys, htable)
		}
		p = e.JsonWrite(p, keys, aline.String())
	}
	return p
}

type Pos struct {
	X, Y           int
	Row, Col, Addr string
}

// Convert 坐标转Excel地址
func (p *Pos) Convert() {
	p.Col = x2col(p.X)
	p.Row = strconv.Itoa(p.Y)
	p.Addr = p.Col + p.Row
}

func (p *Pos) String() string {
	// 先转化
	p.Convert()
	return p.Addr
}

func x2col(x int) string {
	result := ""
	for x > 0 {
		x--
		result = string(rune(x%26+'A')) + result
		x = x / 26
	}
	return result
}

// TODO: excel地址转坐标

func JsonKeys(json string) []string {
	keys := make([]string, 0)
	gjson.Parse(json).ForEach(func(key, value gjson.Result) bool {
		keys = append(keys, key.String())
		return true
	})
	return keys
}
