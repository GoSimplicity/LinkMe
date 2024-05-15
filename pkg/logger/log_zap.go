package logger

import "go.uber.org/zap"

type ZapLogger struct {
	l *zap.Logger
}
type Field struct {
	Key string
	Val any
}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	return &ZapLogger{
		l: l,
	}
}

// Debug 记录 debug 级别的日志
func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.l.Debug(msg, z.toArgs(args)...)
}

// Info 记录 info 级别的日志
func (z *ZapLogger) Info(msg string, args ...Field) {
	z.l.Info(msg, z.toArgs(args)...)
}

// Warn 记录 warn 级别的日志
func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.l.Warn(msg, z.toArgs(args)...)
}

// Error 记录 error 级别的日志
func (z *ZapLogger) Error(msg string, args ...Field) {
	z.l.Error(msg, z.toArgs(args)...)
}

// toArgs 方法将自定义类型 Field 切片转换为 zap.Field 切片
func (z *ZapLogger) toArgs(args []Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Val)) // 将每个 Field 对象转换为 zap.Any 类型的 zap.Field 并添加到结果切片中
	}
	return res
}
