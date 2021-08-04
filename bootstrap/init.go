package bootstrap

import (
	"export-server/models/dao"
	"export-server/models/data"
	"export-server/pkg/conf"
	"export-server/pkg/glog"
)

// InitConf 配置文件初始化
func InitConf(filePath *string) {
	err := conf.InitAppConf(filePath)
	if err != nil {
		panic(err)
	}
}

// initLog 初始化日志
func InitLog() {
	c := conf.AppConf
	if c.GetString("log.type") == "file" {
		glog.InitLog2file(
			c.GetString("log.path"),
			c.GetString("log.level"),
		)
	} else {
		glog.InitLog2std(c.GetString("log.level"))
	}
}

// InitDB 初始化db
func InitDB() {
	dao.MysqlInit()

	err := dao.InitRedis()
	if err != nil {
		panic(err)
	}
}

// InitConsumer 初始化消费者
func InitConsumer() {
	c := conf.AppConf
	new(data.HttpWorker).Run(c.GetInt("taskPool.httpWorker"))
	new(data.RawWorker).Run(c.GetInt("taskPool.rowWorker"))
}
