package global

const (
	TaskHttpKey = "task:http_q" // http的redis任务队列key
	TaskRawKey  = "task:raw_q"  // 请求中携带的源数据的任务队列key
)
