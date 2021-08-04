package bootstrap

import (
	"export-server/models/dao"
	"export-server/models/dao/aoss"
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
	c := conf.AppConf
	dao.MysqlInit()

	err := dao.InitRedis()
	if err != nil {
		panic(err)
	}

	switch c.GetString("exp_storage.channel") {
	case "local":
		// 初始化本地存储层
	case "alioss":
		// 初始化阿里云oss
		aoss.InitAlioss(
			c.GetString("exp_storage.alioss.endpoint_up"),
			c.GetString("exp_storage.alioss.key"),
			c.GetString("exp_storage.alioss.secret"),
			c.GetString("exp_storage.alioss.bucket"),
		)
	default:
		panic("[storage] no such storage")
	}

}

// InitConsumer 初始化消费者
func InitConsumer() {
	c := conf.AppConf
	new(data.HttpWorker).Run(c.GetInt("taskPool.httpWorker"))
	new(data.RawWorker).Run(c.GetInt("taskPool.rowWorker"))
}
