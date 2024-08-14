package ioc

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/job"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

func InitRanking(l *zap.Logger, svc service.RankingService) *cron.Cron {
	// 初始化 cron 实例
	c := cron.New()

	// 定义 TopN 任务逻辑
	topNTask := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		return svc.TopN(ctx)
	}

	jobBuilder := job.NewCronJobBuilder("TopNTask", topNTask, l)
	schedule := "*/30 * * * *" // 每30分钟运行一次任务

	_, err := job.StartCronJob(c, jobBuilder, schedule)
	if err != nil {
		l.Error("Failed to start cron job", zap.Error(err))
		return nil
	}

	return c
}
