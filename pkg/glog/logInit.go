package glog

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	//ZapLoger  *zap.Logger 日志记录器
	ZapLoger *zap.Logger
	//Level 日志最小级别
	Level = zap.NewAtomicLevel()
)

func InitLog2std(level string) {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	SetAtomLevel(level)
	core := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), Level)
	ZapLoger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func InitLog2file(filename, level string) {
	// 按天切割日志
	writer := fileWriterByDay(filename)
	// json格式
	encoder := fileLogEncoder()
	// 动态设置日志级别
	SetAtomLevel(level)

	core := zapcore.NewCore(encoder, writer, Level)
	ZapLoger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// SetLevel 动态设置日志级别 level=[debug,info,warn,error]
func SetAtomLevel(level string) {
	loglevel := zapcore.InfoLevel
	switch strings.ToLower(level) {
	case "debug":
		loglevel = zapcore.DebugLevel
	case "info":
		loglevel = zapcore.InfoLevel
	case "warn":
		loglevel = zapcore.WarnLevel
	case "error":
		loglevel = zapcore.ErrorLevel
	default:
		loglevel = zapcore.InfoLevel
	}
	Level.SetLevel(loglevel)
}

func fileLogEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	// 覆盖默认配置
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(config)
}

func fileWriterByDay(filename string) zapcore.WriteSyncer {
	ext := filepath.Ext(filename)
	path := filepath.Dir(filename)
	file := filepath.Base(filename)
	filebase := file[:len(file)-len(ext)]
	filename = filebase + "-%Y-%m-%d" + ext
	filename = filepath.Join(path, filename)

	hook, err := rotatelogs.New(
		filename,
		rotatelogs.WithMaxAge(time.Hour*24*365),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(hook)
}
