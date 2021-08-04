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
	return "export_file"
}

func (e *ExportFile) DownUrl(key string) string {
	res := dao.MDB.First(e, "hash_key=?", key)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return ""
	}
	if res.Error != nil {
		return ""
	}
	return HttpAbsUrl(e.Path)
}

func HttpAbsUrl(path string) string {
	c := conf.AppConf
	var url string
	// excel 文件存储层
	switch c.GetString("exp_storage.channel") {
	case "local":
		// 初始化本地存储层
		url = c.GetString("exp_storage.local.down_url") + "/v1/down/" + path
	case "alioss":
		// 初始化阿里云oss
		url = c.GetString("exp_storage.alioss.down_url") + "/" + path
	default:
		url = path
	}
	return url
}
