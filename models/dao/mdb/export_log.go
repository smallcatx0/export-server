package mdb

import (
	"errors"
	"export-server/models/dao"
	"time"

	"gorm.io/gorm"
)

type ExportLog struct {
	Id         int       `gorm:"column:id;primary_key;auto_increment" json:"id"`
	HashKey    string    `gorm:"column:hash_key" json:"hash_key"`       //参数哈希
	Title      string    `gorm:"column:title" json:"title"`             //导出标题
	ExtType    string    `gorm:"column:ext_type" json:"ext_type"`       //导出类型(文件后缀)
	SourceType string    `gorm:"column:source_type" json:"source_type"` //数据源类型
	Param      string    `gorm:"column:param" json:"param"`             //请求参数（json）
	UserId     string    `gorm:"column:user_id" json:"user_id"`         //用户id
	Callback   string    `gorm:"column:callback" json:"callback"`       //回调地址
	Status     int       `gorm:"column:status" json:"status"`           //状态：1处理中 2导出成功 3导出失败 4导出取消
	FailReason string    `gorm:"column:fail_reason" json:"fail_reason"` //失败理由
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (e *ExportLog) TableName() string {
	return "export_log"
}

const (
	ExportLog_status_pending = 1 // 处理中
	ExportLog_status_succ    = 2 // 导出成功
	ExportLog_status_fail    = 3 // 导出失败
	ExportLog_status_cancle  = 4 // 导出取消
)

var (
	ExportLog_Stype_Http = "http"
	ExportLog_Stype_Sql  = "sql"
	ExportLog_Stype_Raw  = "raw"
	ExportLog_Stypes     = []string{
		ExportLog_Stype_Http,
		ExportLog_Stype_Sql,
		ExportLog_Stype_Raw,
	}
)

func (e *ExportLog) HashHeyExisted(hashkey string) bool {
	err := dao.MDB.Model(e).Where("hash_key = ?", hashkey).First(e).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// 保存失败理由
func (e *ExportLog) SaveFailReason(reason string) error {
	db := dao.MDB.Model(e).Updates(
		map[string]interface{}{
			"status":      ExportLog_status_fail,
			"fail_reason": reason,
		},
	)
	return db.Error
}
