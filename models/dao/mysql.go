package dao

import (
	"log"
	"time"

	"export-server/pkg/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MDB *gorm.DB

func MysqlInit() {
	c := conf.AppConf
	dsn := c.GetString("mysql.dsn")
	debug := c.GetBool("mysql.debug")
	maxIdleConns := c.GetInt("mysql.maxIdleConns")
	maxOpenConns := c.GetInt("mysql.maxOpenConns")
	connMaxLifetime := c.GetInt("mysql.connMaxLifetime")

	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("[dao] open connDB err ", err)
	}
	if debug {
		Db = Db.Debug()
	}
	sqlDb, err := Db.DB()
	if err != nil {
		log.Panic("[dao] get mysqlDb err ", err)
	}
	sqlDb.SetMaxIdleConns(maxIdleConns)
	sqlDb.SetMaxOpenConns(maxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	err = sqlDb.Ping()
	if err != nil {
		log.Panic("[dao] ping mysql err ", err)
	}
	log.Print("[dao] ping mysql err ")
	MDB = Db
}
