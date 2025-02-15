package es

import (
	"context"
	"database/sql"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"reflect"
	"sync"
	"time"
)

const (
	DEFAULT_FLUSH_TIMEOUT = 1 * time.Second
	DEFAULT_BULK_SIZE     = 1000
)

type FlushUser struct {
	userBuffer   []User        //用户数据刷新区
	bulkSize     int           //批量插入的阈值
	flushTimeout time.Duration //刷新间隔
	bufferMutex  sync.Mutex    //缓冲区互斥锁
	lastFlush    time.Time     //最近刷新时间
	ticker       *time.Ticker  //定时触发器
	done         chan struct{} //关闭信号

	esConsumer *EsConsumer //回调EsConsumer的函数，对消息进行消费
}

type FlushPost struct {
	postBuffer   []Post        //帖子数据刷新区
	bulkSize     int           //批量插入的阈值
	flushTimeout time.Duration //刷新间隔
	bufferMutex  sync.Mutex    //缓冲区互斥锁
	lastFlush    time.Time     //最近刷新时间
	ticker       *time.Ticker  //定时触发器
	done         chan struct{} //关闭信号

	esConsumer *EsConsumer //回调EsConsumer的函数，对消息进行消费
}

func NewFlushUser(esConsumer *EsConsumer) *FlushUser {
	flush := &FlushUser{
		flushTimeout: DEFAULT_FLUSH_TIMEOUT,
		bulkSize:     DEFAULT_BULK_SIZE,
		done:         make(chan struct{}),

		esConsumer: esConsumer,
	}
	flush.ticker = time.NewTicker(flush.flushTimeout)
	go flush.autoFlush()
	return flush
}

func (f *FlushUser) Close() {
	close(f.done)   //发送关闭信号
	f.ticker.Stop() //停止定时器
	f.flushBuffer() //关闭前强制刷新剩余数据
}

// 自动刷新缓冲区的协程
func (f *FlushUser) autoFlush() {
	for {
		select {
		case <-f.ticker.C:
			f.flushBuffer()
		case <-f.done:
			return //退出协程
		}
	}
}

func (f *FlushUser) flushBuffer() {
	f.bufferMutex.Lock()
	defer f.bufferMutex.Unlock()

	if len(f.userBuffer) == 0 {
		return //无数据则直接返回
	}

	if err := f.esConsumer.bulkInsertUser(context.Background(), f.userBuffer); err != nil {
		f.Close()
		f.esConsumer.l.Error("批量插入User失败", zap.Error(err))
	} else {
		f.esConsumer.l.Info("批量插入User成功", zap.Int("count", len(f.userBuffer)))
	}

	f.userBuffer = nil
	f.lastFlush = time.Now()
}

func NewFlushPost(esConsumer *EsConsumer) *FlushPost {
	flush := &FlushPost{
		flushTimeout: DEFAULT_FLUSH_TIMEOUT,
		bulkSize:     DEFAULT_BULK_SIZE,
		done:         make(chan struct{}),

		esConsumer: esConsumer,
	}
	flush.ticker = time.NewTicker(flush.flushTimeout)
	go flush.autoFlush()
	return flush
}

func (f *FlushPost) Close() {
	close(f.done)   //发送关闭信号
	f.ticker.Stop() //停止定时器
	f.flushBuffer() //关闭前强制刷新剩余数据
}

// 自动刷新缓冲区的协程
func (f *FlushPost) autoFlush() {
	for {
		select {
		case <-f.ticker.C:
			f.flushBuffer()
		case <-f.done:
			return //退出协程
		}
	}
}

func (f *FlushPost) flushBuffer() {
	f.bufferMutex.Lock()
	defer f.bufferMutex.Unlock()

	if len(f.postBuffer) == 0 {
		return //无数据则直接返回
	}

	if err := f.esConsumer.bulkInsertPost(context.Background(), f.postBuffer); err != nil {
		f.esConsumer.l.Error("批量插入User失败", zap.Error(err))
	} else {
		f.esConsumer.l.Info("批量插入User成功", zap.Int("count", len(f.postBuffer)))
	}

	f.postBuffer = nil
	f.lastFlush = time.Now()
}

// decodeEventDataToPosts 解析事件数据为 Post 结构体
func decodeEventDataToPost(data interface{}, posts *Post) error {
	config := &mapstructure.DecoderConfig{
		Result:           posts,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToTimeHookFunc("2006-01-02 15:04:05.999"),
			stringToNullTimeHookFunc("2006-01-02 15:04:05.999"),
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(data)
}

// decodeEventDataToUser 解析事件数据为 User 结构体
func decodeEventDataToUser(data interface{}, users *User) error {
	config := &mapstructure.DecoderConfig{
		Result:           users,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToNullTimeHookFunc("2006-01-02 15:04:05.999"),
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(data)
}

// stringToTimeHookFunc 转换字符串到时间类型
func stringToTimeHookFunc(layout string) mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
		if f != reflect.String || t != reflect.Struct {
			return data, nil
		}

		str := data.(string)
		if str == "" {
			return time.Time{}, nil
		}

		return time.Parse(layout, str)
	}
}

// stringToNullTimeHookFunc 转换字符串到 NullTime 类型
func stringToNullTimeHookFunc(layout string) mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
		if f != reflect.String || t != reflect.Struct {
			return data, nil
		}

		str := data.(string)
		if str == "" {
			return sql.NullTime{Valid: false}, nil
		}

		parsedTime, err := time.Parse(layout, str)
		if err != nil {
			return nil, err
		}
		return sql.NullTime{Time: parsedTime, Valid: true}, nil
	}
}
