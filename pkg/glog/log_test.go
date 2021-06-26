package glog_test

import (
	"testing"

	"export-server/pkg/glog"
)

func TestFile(t *testing.T) {
	glog.InitLog2file("/home/logs/tank/curr.log", "Info")
	param := map[string]interface{}{
		"name": "kui",
		"age":  18,
	}

	glog.Debug("该条日志不应会被记录")
	glog.SetAtomLevel("debug")

	glog.Debug("test debug ")
	glog.Debug("test debug with requestId", "requestId")
	glog.Debug("test debug with more", "requestId", "extra one", "extra two")
	glog.DebugT("test debug json", "requestId", param, param)
	glog.DebugF("test debug template age=%d", "requestId", 23)

	glog.Info("测试INFO 级别完整信息", "requestId", "扩展信息1", "扩展信息2")
	glog.InfoF("测试模板日志name=%s", "requestId", "kui")
	glog.InfoT("测试模板日志Json扩展信息", "requestId", param, param)

	glog.Warn("测试warn级别完整信息", "requestId", "扩展信息1", "扩展信息2")
	glog.WarnF("测试模板日志name=%s", "requestId", "kui")
	glog.WarnT("测试模板日志Json扩展信息", "requestId", param, param)

	glog.Error("测试warn级别完整信息", "requestId", "扩展信息1", "扩展信息2")
	glog.ErrorF("测试模板日志name=%s", "requestId", "kui")
	glog.ErrorT("测试模板日志Json扩展信息", "requestId", param, param)

	glog.DPanic("测试DPanic级别完整信息", "request-dadfmwesd", "扩展信息1", "扩展信息2")

	glog.Sync()
}

func TestCons(t *testing.T) {
	defer func() {
		recover()
	}()
	glog.InitLog2std("info")
	glog.Debug("不会被打出来的日志")
	glog.SetAtomLevel("debug")
	glog.Debug("会被打出来的 debug日志", "requestId", "extra1", "extra2")
	glog.Panic("Panic日志", "requestId", "extra1", "extra2")
	t.Log("test ok")
}
