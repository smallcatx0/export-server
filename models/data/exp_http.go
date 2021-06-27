package data

import (
	"export-server/models/dao/rdb"

	request "gitee.com/smallcatx0/gequest"
)

var RdbMq *rdb.Mq
var HttpCli *request.Core

func InitHttpSub(pool int) {

}
