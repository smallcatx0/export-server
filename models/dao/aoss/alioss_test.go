package aoss_test

import (
	"export-server/models/dao"
	"export-server/models/dao/aoss"
	"testing"

	"github.com/golang-module/carbon"
	"github.com/stretchr/testify/assert"
)

var fs dao.FileStorage
var alioss *aoss.AliOssStore

func setup() {
	alioss = aoss.NewAliOssStore("oss-cn-beijing.aliyuncs.com", "LTAIo3aI8hHfCfhX", "VHYxiMjjOTl9VUwdqg6OfBpmQ1RjRO", "heykui")
	fs = alioss
}

func TestUpFile(t *testing.T) {
	ass := assert.New(t)
	setup()
	localFile := "D:\\tmp\\outexcel\\8acbf7ec30e225f570e943e541323d3d.zip"
	err := alioss.Bucket.PutObjectFromFile("export/test.zip", localFile)
	ass.NoError(err)
}

func TestExportPut(t *testing.T) {
	ass := assert.New(t)
	setup()
	localFile := "D:\\tmp\\outexcel\\8acbf7ec30e225f570e943e541323d3d.zip"
	objname, err := fs.Put(localFile)
	ass.Equal(objname, "export/"+carbon.Now().Format("Ym/d/")+"8acbf7ec30e225f570e943e541323d3d.zip")
	ass.NoError(err)
}
