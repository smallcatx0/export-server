package exception

const (
	Succ        = 200
	Fail        = 400
	CalFail     = 609010 // 外部依赖请求失败
	CalShopActs = 609011 // 优惠信息查询失败

)

// ErrNos 错误码映射
var ErrNos = map[uint32]string{
	Succ:   "操作成功",
	Fail:   "系统错误",
	609001: "mysql 连接失败",
	609002: "redis 连接失败",
	609003: "es 连接失败",
	609004: "课程优惠信息查询失败",
}

var (
	ErrMysql  = NewError(609001, 400)
	ErrRedis  = NewError(609002, 400)
	ErrEsPing = NewError(609003, 400)
)
