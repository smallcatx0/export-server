package mdb

import "time"

type ExportLog struct {
	Id         int       `gorm:"column:id"`
	HashKey    string    `gorm:"column:hash_key"`    //参数哈希
	Title      string    `gorm:"column:title"`       //导出标题
	ExtType    string    `gorm:"ext_type"`           // 导出类型
	SourceType string    `gorm:"column:source_type"` //数据源类型
	Param      string    `gorm:"column:param"`       //请求参数（json）
	UserId     string    `gorm:"column:user_id"`     //用户id
	Callback   string    `gorm:"column:callback"`    //回调地址
	CreateAt   time.Time `gorm:"column:create_at"`
	UpdateAt   time.Time `gorm:"column:update_at"`
}

func (e *ExportLog) TableName() string {
	return "export_log"
}
