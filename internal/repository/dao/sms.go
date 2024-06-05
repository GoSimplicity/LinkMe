package dao

import (
	. "LinkMe/internal/repository/models"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SmsDAO interface {
	Insert(ctx context.Context, log VCodeSmsLog) error
	FindFailedLogs(ctx context.Context) []VCodeSmsLog //查找当前时刻以前，发送失败的logs，后续需要重新发送
	Update(ctx context.Context, log VCodeSmsLog) error
}

type smsDao struct {
	db *gorm.DB
	l  *zap.Logger
}

func NewSmsDAO(db *gorm.DB, l *zap.Logger) SmsDAO {
	return &smsDao{
		db: db,
		l:  l,
	}
}

func (s smsDao) Insert(ctx context.Context, log VCodeSmsLog) error {
	//TODO implement me
	panic("implement me")
}

func (s smsDao) FindFailedLogs(ctx context.Context) []VCodeSmsLog {
	//TODO implement me
	panic("implement me")
}

func (s smsDao) Update(ctx context.Context, log VCodeSmsLog) error {
	//TODO implement me
	panic("implement me")
}
