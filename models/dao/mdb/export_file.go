package mdb

import (
	"errors"
	"export-server/models/dao"
	"export-server/pkg/conf"
	"time"

	"gorm.io/gorm"
)

type ExportFile struct {
	Id        int       `gorm:"column:id" json:"id"`
	HashKey   string    `gorm:"column:hash_key" json:"hash_key"`
	Path      string    `gorm:"column:path" json:"path"`
	Type      string    `gorm:"column:type" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"create_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"update_at"`
}

func (e *ExportFile) TableName() string {
	return "base_export_file"
}

func (e *ExportFile) DownUrl(key string) string {
	res := dao.MDB.First(e, "hash_key=?", key)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return ""
	}
	if res.Error != nil {
		return ""
	}
	return OssAbsUrl(e.Path)
}

func OssAbsUrl(path string) string {
	return conf.AppConf.GetString("alioss.excel.endpoint_down") + "/" + path
}
