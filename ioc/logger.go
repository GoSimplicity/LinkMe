package ioc

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() *zap.Logger {
	// 从配置获取日志目录
	logDir := viper.GetString("log.dir")
	logFile := filepath.Join(logDir, "linkme-"+time.Now().Format("2006-01-02")+"-json.log")

	// 日志轮转配置
	fileWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,   // 每个日志文件最大10MB
		MaxBackups: 5,    // 保留最近5个日志文件
		MaxAge:     30,   // 日志文件最多保留30天
		Compress:   true, // 压缩旧日志
		LocalTime:  true, // 使用本地时间
	}

	// 多路输出配置
	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(fileWriter),
	)

	// 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 创建核心配置
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zap.NewAtomicLevelAt(zapcore.InfoLevel),
	)

	// 创建并返回logger
	return zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}
