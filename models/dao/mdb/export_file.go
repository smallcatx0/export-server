package mdb

import "time"

type ExportFile struct {
	Id       int       `gorm:"column:id" json:"id"`
	HashKey  string    `gorm:"column:hash_key" json:"hash_key"`
	Path     string    `gorm:"column:path" json:"path"`
	Type     string    `gorm:"column:type" json:"type"`
	CreateAt time.Time `gorm:"column:create_at" json:"create_at"`
	UpdateAt time.Time `gorm:"column:update_at" json:"update_at"`
}

func (e *ExportFile) TableName() string {
	return "export_file"
}
