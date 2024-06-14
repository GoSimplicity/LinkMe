package dao

import (
	. "LinkMe/internal/repository/models"
	"context"
	"gorm.io/gorm"
	"time"
)

const (
	// jobStatusWaiting 任务处于等待状态
	jobStatusWaiting = iota
	// jobStatusRunning 任务正在运行
	jobStatusRunning
	// jobStatusPaused 任务已暂停
	jobStatusPaused
)

// JobDAO 定义了任务数据访问对象接口
type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, jobId int64) error
	UpdateTime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
}

// jobDAO 实现了 JobDAO 接口
type jobDAO struct {
	db *gorm.DB
}

// NewJobDAO 创建并初始化 jobDAO 实例
func NewJobDAO(db *gorm.DB) JobDAO {
	return &jobDAO{
		db: db,
	}
}

// Preempt 抢占一个等待状态的任务
func (dao *jobDAO) Preempt(ctx context.Context) (Job, error) {
	db := dao.db.WithContext(ctx)
	for {
		var j Job
		now := time.Now().UnixMilli()
		// 查找一个等待状态且下一次执行时间小于当前时间的任务
		err := db.Where("status = ? AND next_time < ?", jobStatusWaiting, now).First(&j).Error
		if err != nil {
			return j, err
		}
		// 尝试更新任务的状态和版本
		result := db.Model(&Job{}).Where("id = ? AND version = ?", j.Id, j.Version).Updates(map[string]any{
			"status":     jobStatusRunning,
			"version":    j.Version + 1,
			"updated_at": now,
		})
		if result.Error != nil {
			return Job{}, result.Error
		}
		if result.RowsAffected == 0 {
			// 如果没有抢到任务，继续循环
			continue
		}
		return j, nil
	}
}

// Release 释放一个正在运行的任务，将其状态重置为等待状态
func (dao *jobDAO) Release(ctx context.Context, jobId int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ?", jobId).Updates(map[string]any{
		"status":     jobStatusWaiting,
		"updated_at": now,
	}).Error
}

// UpdateTime 更新任务的更新时间
func (dao *jobDAO) UpdateTime(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ?", id).Updates(map[string]any{
		"updated_at": now,
	}).Error
}

// UpdateNextTime 更新任务的下次执行时间
func (dao *jobDAO) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ?", id).Updates(map[string]any{
		"updated_at": now,
		"next_time":  t.UnixMilli(),
	}).Error
}
