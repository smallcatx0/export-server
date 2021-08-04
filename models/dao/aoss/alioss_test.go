package aoss_test

import (
	"export-server/models/dao/aoss"
	"testing"

	"github.com/golang-module/carbon"
	"github.com/stretchr/testify/assert"
)

func setup() {
	aoss.InitAlioss("oss-cn-beijing.aliyuncs.com", "", "", "heykui")
}

func TestUpFile(t *testing.T) {
	ass := assert.New(t)
	setup()
	localFile := "D:\\tmp\\outexcel\\9f5ff01641d4ebf41273e9cba7da2dc1.zip"
	err := aoss.Aoss.Bucket.PutObjectFromFile("export/test.zip", localFile)
	ass.NoError(err)
}

func TestExportPut(t *testing.T) {
	ass := assert.New(t)
	setup()
	localFile := "D:\\tmp\\outexcel\\9f5ff01641d4ebf41273e9cba7da2dc1.zip"
	objname, err := aoss.PutExportFile(localFile)
	ass.Equal(objname, "export/"+carbon.Now().Format("Ym/d/")+"9f5ff01641d4ebf41273e9cba7da2dc1.zip")
	ass.NoError(err)
}
