package glog

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

func Debug(msg string, extra ...string) {
	var request_id string
	if len(extra) >= 1 {
		request_id = extra[0]
		extra = extra[1:]
	}
	ZapLoger.Debug(msg, zap.String("request_id", request_id), zap.Strings("extra", extra))
}

func DebugF(template, requestID string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	ZapLoger.Debug(msg, zap.String("request_id", requestID))
}

func DebugT(msg, requestID string, extra ...interface{}) {
	extSlice := make([]string, 0, len(extra))
	for _, one := range extra {
		tmpStr, _ := json.Marshal(one)
		extSlice = append(extSlice, string(tmpStr))
	}
	ZapLoger.Debug(msg, zap.String("request_id", requestID), zap.Strings("extra", extSlice))
}

func Info(msg string, extra ...string) {
	var request_id string
	if len(extra) >= 1 {
		request_id = extra[0]
		extra = extra[1:]
	}
	ZapLoger.Info(msg, zap.String("request_id", request_id), zap.Strings("extra", extra))
}

func InfoF(template, requestID string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	ZapLoger.Info(msg, zap.String("request_id", requestID))
}

func InfoT(msg, requestID string, extra ...interface{}) {
	extSlice := make([]string, 0, len(extra))
	for _, one := range extra {
		tmpStr, _ := json.Marshal(one)
		extSlice = append(extSlice, string(tmpStr))
	}
	ZapLoger.Info(msg, zap.String("request_id", requestID), zap.Strings("extra", extSlice))
}

func Warn(msg string, extra ...string) {
	var request_id string
	if len(extra) >= 1 {
		request_id = extra[0]
		extra = extra[1:]
	}
	ZapLoger.Warn(msg, zap.String("request_id", request_id), zap.Strings("extra", extra))
}

func WarnF(template, requestID string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	ZapLoger.Warn(msg, zap.String("request_id", requestID))
}

func WarnT(msg, requestID string, extra ...interface{}) {
	extSlice := make([]string, 0, len(extra))
	for _, one := range extra {
		tmpStr, _ := json.Marshal(one)
		extSlice = append(extSlice, string(tmpStr))
	}
	ZapLoger.Warn(msg, zap.String("request_id", requestID), zap.Strings("extra", extSlice))
}

func Error(msg string, extra ...string) {
	var request_id string
	if len(extra) >= 1 {
		request_id = extra[0]
		extra = extra[1:]
	}
	ZapLoger.Error(msg, zap.String("request_id", request_id), zap.Strings("extra", extra))
}

func ErrorF(template, requestID string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	ZapLoger.Error(msg, zap.String("request_id", requestID))
}

func ErrorT(msg, requestID string, extra ...interface{}) {
	extSlice := make([]string, 0, len(extra))
	for _, one := range extra {
		tmpStr, _ := json.Marshal(one)
		extSlice = append(extSlice, string(tmpStr))
	}
	ZapLoger.Error(msg, zap.String("request_id", requestID), zap.Strings("extra", extSlice))
}

func DPanic(msg string, extra ...string) {
	var request_id string
	if len(extra) >= 1 {
		request_id = extra[0]
		extra = extra[1:]
	}
	ZapLoger.DPanic(msg, zap.String("request_id", request_id), zap.Strings("extra", extra))
}

func DPanicF(template, requestID string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	ZapLoger.DPanic(msg, zap.String("request_id", requestID))
}

func DPanicT(msg, requestID string, extra ...interface{}) {
	extSlice := make([]string, 0, len(extra))
	for _, one := range extra {
		tmpStr, _ := json.Marshal(one)
		extSlice = append(extSlice, string(tmpStr))
	}
	ZapLoger.DPanic(msg, zap.String("request_id", requestID), zap.Strings("extra", extSlice))
}

func Panic(msg string, extra ...string) {
	var request_id string
	if len(extra) >= 1 {
		request_id = extra[0]
		extra = extra[1:]
	}
	ZapLoger.Panic(msg, zap.String("request_id", request_id), zap.Strings("extra", extra))
}

func PanicF(template, requestID string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	ZapLoger.Panic(msg, zap.String("request_id", requestID))
}

func PanicT(msg, requestID string, extra ...interface{}) {
	extSlice := make([]string, 0, len(extra))
	for _, one := range extra {
		tmpStr, _ := json.Marshal(one)
		extSlice = append(extSlice, string(tmpStr))
	}
	ZapLoger.Panic(msg, zap.String("request_id", requestID), zap.Strings("extra", extSlice))
}

func Fatal(msg string, extra ...string) {
	var request_id string
	if len(extra) >= 1 {
		request_id = extra[0]
		extra = extra[1:]
	}
	ZapLoger.Fatal(msg, zap.String("request_id", request_id), zap.Strings("extra", extra))
}

func FatalF(template, requestID string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	ZapLoger.Fatal(msg, zap.String("request_id", requestID))
}

func FatalT(msg, requestID string, extra ...interface{}) {
	extSlice := make([]string, 0, len(extra))
	for _, one := range extra {
		tmpStr, _ := json.Marshal(one)
		extSlice = append(extSlice, string(tmpStr))
	}
	ZapLoger.Fatal(msg, zap.String("request_id", requestID), zap.Strings("extra", extSlice))
}

func Sync() {
	ZapLoger.Sync()
}
