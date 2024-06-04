package dao

import (
	. "LinkMe/internal/repository/models"
	"context"
	"gorm.io/gorm"
)

type SendVCodeDAO interface {
	Insert(ctx context.Context, log VCodeSmsLog) error
	FindFailedLogs(ctx context.Context) []VCodeSmsLog //查找当前时刻以前，发送失败的logs，后续需要重新发送
	Update(ctx context.Context, log VCodeSmsLog)
}

type sendVCodeDAO struct {
	db *gorm.DB
}

func NewGORMAsyncSmsDAO(db *gorm.DB) SendVCodeDAO {
	return &sendVCodeDAO{
		db: db,
	}
}

func (s sendVCodeDAO) Insert(ctx context.Context, log VCodeSmsLog) error {
	//TODO implement me
	panic("implement me")
}

func (s sendVCodeDAO) FindFailedLogs(ctx context.Context) []VCodeSmsLog {
	//TODO implement me
	panic("implement me")
}

func (s sendVCodeDAO) Update(ctx context.Context, log VCodeSmsLog) {
	//TODO implement me
	panic("implement me")
}
