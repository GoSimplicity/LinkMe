/*
 * MIT License
 *
 * Copyright (c) 2024 Bamboo
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

package di

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化日志系统并返回一个配置好的logger
func InitLogger() *zap.Logger {
	// 从配置获取日志目录，如果未配置则使用默认值
	logDir := viper.GetString("log.dir")
	if logDir == "" {
		logDir = "logs"
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("创建日志目录失败: " + err.Error())
	}

	logFile := filepath.Join(logDir, "linkme-"+time.Now().Format("2006-01-02")+".log")
	logLevel := getLogLevel(viper.GetString("log.level"))

	// 日志轮转配置
	fileWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    viper.GetInt("log.max_size"),    // 每个日志文件最大尺寸，默认10MB
		MaxBackups: viper.GetInt("log.max_backups"), // 保留的旧日志文件数量
		MaxAge:     viper.GetInt("log.max_age"),     // 日志文件保留天数
		Compress:   viper.GetBool("log.compress"),   // 是否压缩旧日志
		LocalTime:  true,                            // 使用本地时间
	}

	if fileWriter.MaxSize == 0 {
		fileWriter.MaxSize = 10 // 默认10MB
	}
	if fileWriter.MaxBackups == 0 {
		fileWriter.MaxBackups = 5 // 默认保留5个备份
	}
	if fileWriter.MaxAge == 0 {
		fileWriter.MaxAge = 30 // 默认保留30天
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
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeDuration = zapcore.MillisDurationEncoder

	// 创建核心配置
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		logLevel,
	)

	// 创建并返回logger
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCallerSkip(1),
	)

	// 替换全局logger
	zap.ReplaceGlobals(logger)

	return logger
}

// getLogLevel 根据配置字符串返回对应的日志级别
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
