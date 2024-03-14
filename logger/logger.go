package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger = initLogger()
var ZapLogger *zap.Logger

// func getFileLogWriter() (writeSyncer zapcore.WriteSyncer) {
// 	// 使用 lumberjack 实现 log rotate
// 	lumberJackLogger := &lumberjack.Logger{
// 		Filename:   "/tmp/test.log",
// 		MaxSize:    100, // 单个文件最大100M
// 		MaxBackups: 60,  // 多于 60 个日志文件后，清理较旧的日志
// 		MaxAge:     1,   // 一天一切割
// 		Compress:   false,
// 	}

// 	return zapcore.AddSync(lumberJackLogger)
// }

// initLogger 初始化zap日志
func initLogger() *zap.SugaredLogger {
	zap.NewDevelopmentConfig()
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.FunctionKey = "Func"

	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
	ZapLogger = logger
	return logger.Sugar()
}
