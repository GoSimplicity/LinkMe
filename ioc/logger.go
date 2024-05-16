package ioc

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogger() *zap.Logger {
	filepath := viper.GetString("log.filepath")
	if filepath == "" {
		fmt.Println("没有找到文件路径")
	}
	// 创建生产环境的编码器配置
	c := zap.NewProductionEncoderConfig()
	c.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(c)
	defaultLogLevel := zapcore.DebugLevel // 默认级别为Debug
	// 创建文件写入器
	logFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		fmt.Printf("无法打开指定路径下的文件: %v", err)
	}
	// 将文件写入器添加到写入器列表中
	writer := zapcore.AddSync(logFile)
	l := zap.New(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zap.AddCaller(),
		// 在ERROR级别添加堆栈跟踪
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	// 后续加入定时任务需定时处理超过指定大小的日志文件

	return l
}
