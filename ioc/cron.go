package ioc

import (
	"github.com/go-co-op/gocron"
	"time"
)

// InitScheduler 初始化并返回一个gocron调度器
func InitScheduler() *gocron.Scheduler {
	// 创建一个新的调度器实例，使用本地时区
	loc, _ := time.LoadLocation("Local")
	scheduler := gocron.NewScheduler(loc)
	return scheduler
}
