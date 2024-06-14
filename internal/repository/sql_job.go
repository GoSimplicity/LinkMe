package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"context"
	"time"
)

// CronJobRepository 定义了定时任务仓库接口
type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, jobId int64) error
	UpdateTime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, nextTime time.Time) error
}

// cronJobRepository 实现了 CronJobRepository 接口
type cronJobRepository struct {
	dao dao.JobDAO
}

// NewCronJobRepository 创建并初始化 cronJobRepository 实例
func NewCronJobRepository(dao dao.JobDAO) CronJobRepository {
	return &cronJobRepository{
		dao: dao,
	}
}

// Preempt 从数据库中抢占一个任务
func (c *cronJobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := c.dao.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	return domain.Job{
		Id:         j.Id,
		Expression: j.Expression,
		Executor:   j.Executor,
		Name:       j.Name,
	}, nil
}

// Release 释放一个任务
func (c *cronJobRepository) Release(ctx context.Context, jobId int64) error {
	return c.dao.Release(ctx, jobId)
}

// UpdateTime 更新任务的时间
func (c *cronJobRepository) UpdateTime(ctx context.Context, id int64) error {
	return c.dao.UpdateTime(ctx, id)
}

// UpdateNextTime 更新任务的下次执行时间
func (c *cronJobRepository) UpdateNextTime(ctx context.Context, id int64, nextTime time.Time) error {
	return c.dao.UpdateNextTime(ctx, id, nextTime)
}
