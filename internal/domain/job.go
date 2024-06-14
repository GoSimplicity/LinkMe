package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	Id         int64  // 任务的唯一标识符
	Name       string // 任务名称
	Expression string // Cron 表达式，用于定义任务的调度时间
	Executor   string // 执行任务的执行器名称
	Cfg        string // 任务配置，可以是任意字符串
	CancelFunc func() // 用于取消任务的函数
}

// NextTime 计算任务的下次执行时间
func (j *Job) NextTime() (time.Time, error) {
	// 创建新的 Cron 表达式解析器
	c := cron.NewParser(cron.Second | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	// 解析 Cron 表达式
	schedule, err := c.Parse(j.Expression)
	if err != nil {
		return time.Time{}, err // 返回解析错误
	}
	// 计算并返回下次执行时间
	return schedule.Next(time.Now()), nil
}
