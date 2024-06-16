package ioc

import (
	"LinkMe/internal/service"
	"LinkMe/job"
	"context"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

func InitRanking(l *zap.Logger, svc service.RankingService) *cron.Cron {
	// 初始化 cron 实例
	c := cron.New(cron.WithSeconds())
	// 定义 TopN 任务逻辑
	topNTask := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		return svc.TopN(ctx)
	}
	// 创建 CronJobBuilder 实例
	jobBuilder := job.NewCronJobBuilder("TopNTask", topNTask, l)
	// 将任务添加到 cron 实例
	schedule := "* * * * *" // 每分钟运行一次任务
	_, err := job.StartCronJob(c, jobBuilder, schedule)
	if err != nil {
		l.Error("Failed to start cron job", zap.Error(err))
		return nil
	}
	// 返回 cron 实例
	return c
}
