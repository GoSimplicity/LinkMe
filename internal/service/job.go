package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
	"go.uber.org/zap"
	"time"
)

type CronJobService interface {
	Preempt(ctx context.Context) (domain.Job, error)
	ResetNextTime(ctx context.Context, dj domain.Job) error
}

type cronJobService struct {
	repo            repository.CronJobRepository
	l               *zap.Logger
	refreshInterval time.Duration
}

func NewCronJobService(repo repository.CronJobRepository, l *zap.Logger) CronJobService {
	return &cronJobService{
		repo:            repo,
		l:               l,
		refreshInterval: time.Minute,
	}
}

// Preempt 获取一个任务，并启动一个协程定期刷新任务的更新时间
func (c *cronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := c.repo.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}

	ticker := time.NewTicker(c.refreshInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.refresh(j.Id)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
	j.CancelFunc = func() {
		ticker.Stop()
		ct, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := c.repo.Release(ct, j.Id)
		if er != nil {
			c.l.Error("Failed to release job", zap.Error(er))
		}
	}
	return j, nil
}

// ResetNextTime 重置任务的下次执行时间
func (c *cronJobService) ResetNextTime(ctx context.Context, dj domain.Job) error {
	nextTime, err := dj.NextTime()
	if err != nil {
		return err
	}
	return c.repo.UpdateNextTime(ctx, dj.Id, nextTime)
}

// refresh 更新任务的更新时间
func (c *cronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := c.repo.UpdateTime(ctx, id)
	if err != nil {
		c.l.Error("Failed to refresh job", zap.Error(err))
	}
}
