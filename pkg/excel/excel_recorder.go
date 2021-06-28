package excel

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tidwall/gjson"
)

type ExcelRecorder struct {
	FilePath string
	fp       *excelize.File
}

func NewExcelRecorder(path string) *ExcelRecorder {
	instance := &ExcelRecorder{
		FilePath: path,
		fp:       excelize.NewFile(),
	}
	return instance
}

func (e *ExcelRecorder) WriteLines(startCell string, lines []map[string]interface{}) {
	// for _, line := range lines {

	// }
	// e.fp.SetCellValue()
}

func Map2Arr(map[string]interface{}) []interface{} {
	// TODO: 实现map2Array
	return nil
}

func X2column(x int) string {
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
