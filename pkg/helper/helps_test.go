package helper_test

import (
	"export-server/pkg/helper"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	assert := assert.New(t)
	assert.True(helper.Empty(0.0))
	assert.True(helper.Empty(0))
	assert.True(helper.Empty(""))
	assert.True(helper.Empty(nil))
	assert.True(helper.Empty(map[string]interface{}{}))
	assert.True(helper.Empty([]string{}))
	assert.False(helper.Empty(' '))
}

func TestDef(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, helper.GetDefInt(1, 10))
	assert.Equal(10, helper.GetDefInt(0, 10))
	assert.Equal("hello", helper.GetDefStr("hello", "world"))
	assert.Equal("world", helper.GetDefStr("", "world"))
}

func TestMap2Arr(t *testing.T) {
	assert := assert.New(t)
	aMap := map[string]interface{}{"Name": "demo", "ID": 1, "Age": 31}
	keys := []string{"ID", "Name", "Age"}
	expect := []interface{}{1, "demo", 31}
	res := helper.Map2Arr(aMap, keys)
	assert.Equal(res, expect)
}

func TestZip(t *testing.T) {
	ssrc := `D:\tmp\outExcel\a807a70cc66eeb3c0ff24450c6caed88`
	helper.FolderZip(ssrc, `D:\tmp\outExcel\out.zip`)
}

func TestEq(t *testing.T) {
	Convey("[fun] EqualInt", t, func() {
		var x interface{} = 1
		So(true, ShouldEqual, helper.EqualInt(x, 1))
		So(false, ShouldEqual, helper.EqualInt(x, 2))

		var s interface{} = "hanmeimei"
		So(true, ShouldEqual, helper.EqualStr(s, "hanmeimei"))
		So(false, ShouldEqual, helper.EqualStr(s, "lilei"))
		So(true, ShouldNotEqual, helper.EqualStr(s, "lilei"))
	})
}
