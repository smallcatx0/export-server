package aoss_test

import (
	"export-server/models/dao"
	"export-server/models/dao/aoss"
	"log"
	"testing"

	"github.com/golang-module/carbon"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var fs dao.FileStorage
var alioss *aoss.AliOssStore
var C *viper.Viper

func ConfInit() {
	C = viper.New()
	C.SetConfigFile("./unittest_env")
	C.SetConfigType("toml")
	err := C.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	if C.GetString("v") != "1.0.0" {
		log.Fatal("读取配置文件失败")
	}
}

func setup() {
	ConfInit()
	alioss = aoss.NewAliOssStore(
		"oss-cn-beijing.aliyuncs.com",
		C.GetString("AOSS_KEY"),
		C.GetString("AOSS_SECRET"),
		"heykui",
	)
	fs = alioss
}

func TestUpFile(t *testing.T) {
	setup()
	ass := assert.New(t)
	localFile := "D:\\tmp\\1.png"
	err := alioss.Bucket.PutObjectFromFile("export/1.png", localFile)
	ass.NoError(err)
}

func TestExportPut(t *testing.T) {
	setup()
	ass := assert.New(t)
	localFile := "D:\\tmp\\1.png"
	objname, err := fs.Put(localFile)
	ass.Equal(objname, "export/"+carbon.Now().Format("Ym/d/")+"1.png")
	ass.NoError(err)
}
