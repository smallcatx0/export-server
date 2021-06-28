package excel_test

import (
	"export-server/pkg/excel"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestX2column(t *testing.T) {
	assert := assert.New(t)
	columnMap := map[int]string{
		1: "A", 2: "B", 9: "I", 32: "AF", 35: "AI", 41: "AO", 61: "BI", 73: "BU",
	}
	for k, act := range columnMap {
		v := excel.X2column(k)
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
