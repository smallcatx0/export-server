package dao_test

import (
	"export-server/models/dao"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalFS(t *testing.T) {
	ass := assert.New(t)
	var Fs dao.FileStorage = &dao.LocalStore{"/store"}
	orgFile := "D:\\tmp\\logs\\file.log"
	objName, err := Fs.Put(orgFile)
	ass.NoError(err)
	log.Print(objName)
	err = os.RemoveAll(orgFile)
	ass.NoError(err)
}
