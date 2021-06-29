package excel

import (
	"encoding/json"
	"export-server/pkg/helper"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tidwall/gjson"
)

type ExcelRecorder struct {
	FilePath string
	ExcelFp  *excelize.File
	Sheet    string
}
type Pos struct {
	X, Y int
	Addr string
}

// TODO: excel地址转坐标

// Convert 坐标转Excel地址
func (p *Pos) Convert() {
	p.Addr = X2col(p.X) + strconv.Itoa(p.Y)
}

func (p *Pos) String() string {
	// 先转化
	p.Convert()
	return p.Addr
}

func NewExcelRecorder(path string) *ExcelRecorder {
	instance := &ExcelRecorder{
		FilePath: path,
		ExcelFp:  excelize.NewFile(),
		Sheet:    "Sheet1",
	}
	return instance
}

func (e *ExcelRecorder) JsonListWrite(p Pos, jsonStr string, isFirst bool) (Pos, error) {
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
		lintMap := line.Value().(map[string]interface{})
		lineValues := helper.Map2Arr(lintMap, keys)
		e.ExcelFp.SetSheetRow(e.Sheet, p.String(), &lineValues)
		p.Y += 1
	}
	return p, nil
}

// JsonListWrite2 效率还不如 JsonListWrite
func (e *ExcelRecorder) JsonListWrite2(p Pos, jsonStr string, isFirst bool) (Pos, error) {
	firstline := gjson.Get(jsonStr, "0").String()
	keys := JsonKeys(firstline)
	lists := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(jsonStr), &lists)
	if err != nil {
		return p, err
	}
	if !isFirst {
		lists = lists[1:]
	}
	mline := List2Arrs(lists, keys)
	e.WriteLines(p, mline)
	return p, nil
}

// WriteLines 多行写入
func (e *ExcelRecorder) WriteLines(p Pos, lines [][]interface{}) Pos {
	for _, line := range lines {
		e.ExcelFp.SetSheetRow(e.Sheet, p.String(), &line)
		p.Y += 1
	}
	return p
}

// Save 保存
func (e *ExcelRecorder) Save() error {
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

func X2col(x int) string {
	result := ""
	for x > 0 {
		x--
		result = string(rune(x%26+'A')) + result
		x = x / 26
	}
	return result
}

func JsonKeys(json string) []string {
	keys := make([]string, 0)
	gjson.Parse(json).ForEach(func(key, value gjson.Result) bool {
		keys = append(keys, key.String())
		return true
	})
	return keys
}
