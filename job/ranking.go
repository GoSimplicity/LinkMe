package job

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type CronJobBuilder interface {
	Name() string
	Run() error
	Build() cron.Job
}

type cronJobBuilder struct {
	name string
	run  func() error
	l    *zap.Logger
}

func NewCronJobBuilder(run func() error, l *zap.Logger) CronJobBuilder {
	return &cronJobBuilder{
		name: "hotSearch",
		run:  run,
		l:    l,
	}
}

// Name 返回任务名称
func (b *cronJobBuilder) Name() string {
	return b.name
}

// Run 执行任务
func (b *cronJobBuilder) Run() error {
	return b.run()
}

// Build 构建 cron.Job
func (b *cronJobBuilder) Build() cron.Job {
	return cron.FuncJob(func() {
		b.l.Debug("开始运行", zap.String("name", b.Name()))
		err := b.Run()
		if err != nil {
			b.l.Error("执行失败", zap.String("name", b.Name()), zap.Error(err))
		}
		b.l.Debug("结束运行", zap.String("name", b.Name()))
	})
}
