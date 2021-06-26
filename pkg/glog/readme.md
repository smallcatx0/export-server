# 日志包





### 文件日志

基于 `go.uber.org/zap`封装的日志包

#### 初始化

```go
// 初始化日志文件
glog.InitLog2file("/home/logs/tank/curr.log", "Info")
// 初始化控制台日志
glog.InitLog2std("Info")

```

#### 使用

```go
// Debug Info Warn Error DPanic Painc Fatal 等级的日志均有如下方法
glog.Debug("msg")
glog.Debug("msg", "requestId", "extra one", "extra two")
glog.DebugT("msg", "requestId", param, param) // param都会被json序列化
glog.DebugF("测试模板日志age=%d", "requestId", 23) 
```

`/home/logs/tank/curr.log`  日志文件中 每行json

```log
{"level":"debug","ts":"2021-05-17T15:09:24.717+0800","caller":"glog/log_test.go:22","msg":"test debug ","request_id":"","extra":[]}
{"level":"debug","ts":"2021-05-17T15:09:24.745+0800","caller":"glog/log_test.go:23","msg":"test debug with requestId","request_id":"requestId","extra":[]}
{"level":"debug","ts":"2021-05-17T15:09:24.745+0800","caller":"glog/log_test.go:24","msg":"test debug with more","request_id":"requestId","extra":["extra one","extra two"]}
{"level":"debug","ts":"2021-05-17T15:09:24.746+0800","caller":"glog/log_test.go:25","msg":"test debug json","request_id":"requestId","extra":["{\"age\":18,\"name\":\"kui\"}","{\"age\":18,\"name\":\"kui\"}"]}
{"level":"debug","ts":"2021-05-17T15:09:24.746+0800","caller":"glog/log_test.go:26","msg":"test debug template age=23","request_id":"requestId"}
...

```



### 钉钉群机器人

#### 初始化

```go
ala := glog.DingAlarmNew(webHook, secret)
```

#### 使用

```go
// 普通消息
ala.Text("测试普通消息").AtPhones("18681636749").Send()
// markdown 消息
ala.Markdown("title", "### 三级标题 \n\n> 引用 \n\n内容").Send()
```

