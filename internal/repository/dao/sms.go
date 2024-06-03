package dao

import (
	. "LinkMe/internal/repository/models"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

const (
	asyncStatusWaiting = iota
	asyncStatusFailed
	asyncStatusSuccess
)

// ErrWaitingSMSNotFound 错误常量定义
var ErrWaitingSMSNotFound = gorm.ErrRecordNotFound

// SendCodeDAO 接口定义了异步 SMS 数据操作的相关方法
type SendCodeDAO interface {
	Insert(ctx context.Context, s Sms) error
	GetWaitingSMS(ctx context.Context) (Sms, error)
	MarkSuccess(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64) error
}

// sendCodeDAO 是使用 GORM 实现的 AsyncSmsDAO
type sendCodeDAO struct {
	db *gorm.DB
}

// NewGORMAsyncSmsDAO 创建并返回一个新的 GORMAsyncSmsDAO 实例
func NewGORMAsyncSmsDAO(db *gorm.DB) SendCodeDAO {
	return &sendCodeDAO{
		db: db,
	}
}

// Insert 插入一个新的异步 SMS 记录
func (g *sendCodeDAO) Insert(ctx context.Context, s Sms) error {
	return g.db.WithContext(ctx).Create(&s).Error
}

// GetWaitingSMS 获取一个待处理的异步 SMS 记录
func (g *sendCodeDAO) GetWaitingSMS(ctx context.Context) (Sms, error) {
	var s Sms
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 避免偶发性的失败，只查找 1 分钟前的记录
		now := time.Now().UnixMilli()
		endTime := now - time.Minute.Milliseconds()
		// 锁定记录，避免并发问题
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("utime < ? AND status = ?", endTime, asyncStatusWaiting).
			First(&s).Error; err != nil {
			return err
		}
		// 更新记录的重试次数和更新时间
		return tx.Model(&Sms{}).
			Where("id = ?", s.Id).
			Updates(map[string]any{
				"retry_cnt": gorm.Expr("retry_cnt + 1"),
				"utime":     now,
			}).Error
	})
	return s, err
}

// MarkSuccess 标记异步 SMS 记录为成功
func (g *sendCodeDAO) MarkSuccess(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&Sms{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"utime":  now,
			"status": asyncStatusSuccess,
		}).Error
}

// MarkFailed 标记异步 SMS 记录为失败
func (g *sendCodeDAO) MarkFailed(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&Sms{}).
		Where("id = ? AND retry_cnt >= retry_max", id).
		Updates(map[string]any{
			"utime":  now,
			"status": asyncStatusFailed,
		}).Error
}
