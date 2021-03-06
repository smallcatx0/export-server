package dao

import (
	"export-server/pkg/helper"
	"os"
	"path/filepath"

	"github.com/golang-module/carbon"
)

var FS FileStorage

type FileStorage interface {
	Put(string) (string, error)
}

// 本地存储器
type LocalStore struct {
	Path string
}

func InitLocaStore(path string) {
	FS = &LocalStore{Path: path}
}

func (f *LocalStore) Put(filename string) (objname string, err error) {
	t := carbon.Now()
	_, name := filepath.Split(filename)
	objname = filepath.Join(
		f.Path,
		t.Format("ym"),
		t.Format("d"),
		name,
	)
	helper.TouchDir(objname)
	err = os.Rename(filename, objname)
	if err != nil {
		return
	}
	objname, err = filepath.Rel(f.Path, objname)
	return
}
