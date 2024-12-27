package ioc

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultBufferSize = 10
	TopicZapLogs      = "linkme_elk_events"
	logLevel          = zap.DebugLevel
)

//// InitLogger 将日志输出到控制台
//func InitLogger(kafkaProducer sarama.SyncProducer) *zap.Logger {
//	kfkcore := NewKafkaCore(TopicZapLogs, kafkaProducer, logLevel)
//	defer func(kfkcore *KafkaCore) {
//		_ = kfkcore.Sync()
//	}(kfkcore)
//	cfg := zap.NewDevelopmentEncoderConfig()
//	encoder := zapcore.NewJSONEncoder(cfg)
//
//	// 创建zap, 在控制台输出的同时将日志发送到kafka中
//	return zap.New(zapcore.NewTee(kfkcore, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel)))
//}
//
//// ----------------kafkaCore 是构造zap的core, 用来消费zap的logs记录--------------------
//
//type KafkaCore struct {
//	producer  sarama.SyncProducer
//	topic     string
//	level     zapcore.Level
//	msgBuffer []*sarama.ProducerMessage
//	maxBuffer int
//	lock      sync.Mutex
//}
//
//type ReadEvent struct {
//	Timestamp int64  `json:"timestamp"`
//	Level     string `json:"level"`
//	Message   string `json:"message"`
//}
//
//func NewKafkaCore(topic string, producer sarama.SyncProducer, level zapcore.Level, bufsize ...int) *KafkaCore {
//	core := &KafkaCore{
//		producer: producer,
//		topic:    topic,
//		level:    level,
//	}
//	if len(bufsize) > 0 {
//		core.maxBuffer = bufsize[0]
//	} else {
//		core.maxBuffer = defaultBufferSize
//	}
//	return core
//}
//
//// Enabled 判断是否达到所指定日志记录的等级
//func (kc *KafkaCore) Enabled(level zapcore.Level) bool {
//	return kc.level.Enabled(level)
//}
//
//// With 该函数主要为Core设置额外的字段,并产生副本
//func (kc *KafkaCore) With(fields []zapcore.Field) zapcore.Core {
//	return &KafkaCore{
//		producer: kc.producer,
//		topic:    kc.topic,
//	}
//}
//
//// Check 该函数用于检查entry是否达到日志记录的等级
//func (kc *KafkaCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
//	if kc.Enabled(entry.Level) {
//		return checked.AddCore(entry, kc)
//	}
//	return checked
//}
//
//// Write 该函数将日志数据输出到kafka消息队列中
//func (kc *KafkaCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
//
//	msg := ReadEvent{
//		Timestamp: entry.Time.Unix(),
//		Level:     entry.Level.String(),
//		Message:   entry.Message,
//	}
//
//	bytes, err := json.Marshal(msg)
//	if err != nil {
//		return err
//	}
//	kafkaMsg := &sarama.ProducerMessage{
//		Topic: kc.topic,
//		Value: sarama.ByteEncoder(bytes),
//	}
//	kc.lock.Lock()
//	defer kc.lock.Unlock()
//	kc.msgBuffer = append(kc.msgBuffer, kafkaMsg)
//	if len(kc.msgBuffer) >= kc.maxBuffer {
//		return kc.flushBuffer()
//	}
//	return nil
//}
//
//// Sync 该函数主要用于将缓存中的日志数据刷新到输出设备（如文件、网络等），确保日志数据被持久化存储或者发送出去
//func (kc *KafkaCore) Sync() error {
//	kc.lock.Lock()
//	defer kc.lock.Unlock()
//	return kc.flushBuffer()
//}
//
//// 将缓冲区的消息批量发送到kafka中
//func (kc *KafkaCore) flushBuffer() error {
//	if len(kc.msgBuffer) == 0 {
//		return nil
//	}
//	var errs []error
//	for _, msg := range kc.msgBuffer {
//		_, _, err := kc.producer.SendMessage(msg)
//		if err != nil {
//			errs = append(errs, err)
//		}
//	}
//	kc.msgBuffer = make([]*sarama.ProducerMessage, 0)
//	if len(errs) > 0 {
//		return fmt.Errorf("failed to send logs to kafka: %v", errs)
//	}
//	return nil
//}

func InitLogger() *zap.Logger {
	// 使用 Lumberjack 进行日志文件滚动
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/linkme.log", // 指定日志文件路径
		MaxSize:    50,                // 每个日志文件的最大大小，单位：MB
		MaxBackups: 3,                 // 保留旧日志文件的最大个数
		MaxAge:     28,                // 保留旧日志文件的最大天数
		Compress:   true,              // 是否压缩旧的日志文件
	}

	// 配置 zap 的日志编码器，格式为 JSON
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"                   // 设置时间字段名为 timestamp
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 设置时间格式为 ISO8601

	// 创建 zap 核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 使用自定义的 JSON 编码器配置
		zapcore.AddSync(lumberjackLogger),     // 日志输出到 Lumberjack 管理的日志文件
		zapcore.DebugLevel,                    // 设置日志级别为 Debug 及以上
	)

	// 创建 Logger 并添加调用者信息
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // 跳过调用栈的第一层，使得调用文件名和行号更准确

	return logger
}
