package ioc

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 将日志输出到控制台
func InitLogger() *zap.Logger {
	// 使用NewDevelopmentConfig创建一个适合开发环境的日志记录器
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 使用彩色输出
	l, _ := cfg.Build()
	return l
}

//func InitLogger() *zap.Logger {
//	// 使用 Lumberjack 进行日志文件滚动
//	lumberjackLogger := &lumberjack.Logger{
//		Filename:   "/var/log/linkme.log", // 指定日志文件路径
//		MaxSize:    50,                    // 每个日志文件的最大大小，单位：MB
//		MaxBackups: 3,                     // 保留旧日志文件的最大个数
//		MaxAge:     28,                    // 保留旧日志文件的最大天数
//		Compress:   true,                  // 是否压缩旧的日志文件
//	}
//
//	// 配置 zap 的日志编码器，格式为 JSON
//	encoderConfig := zap.NewProductionEncoderConfig()
//	encoderConfig.TimeKey = "timestamp"                   // 设置时间字段名为 timestamp
//	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 设置时间格式为 ISO8601
//
//	// 创建 zap 核心
//	core := zapcore.NewCore(
//		zapcore.NewJSONEncoder(encoderConfig), // 使用自定义的 JSON 编码器配置
//		zapcore.AddSync(lumberjackLogger),     // 日志输出到 Lumberjack 管理的日志文件
//		zapcore.DebugLevel,                    // 设置日志级别为 Debug 及以上
//	)
//
//	// 创建 Logger 并添加调用者信息
//	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // 跳过调用栈的第一层，使得调用文件名和行号更准确
//
//	return logger
//}
