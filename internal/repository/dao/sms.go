package dao

import (
	. "LinkMe/internal/repository/models"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type SmsDAO interface {
	Insert(ctx context.Context, log VCodeSmsLog) error
	FindFailedLogs(ctx context.Context) ([]VCodeSmsLog, error) //查找当前时刻以前，发送失败的logs，后续需要重新发送
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

func (s *smsDao) Insert(ctx context.Context, log VCodeSmsLog) error {
	log.UpdatedTime = time.Now().Unix() //初始化插入时的时间戳
	log.CreateTime = time.Now().Unix()  //初始化插入时的时间戳
	return s.db.WithContext(ctx).Create(&log).Error
}

func (s *smsDao) FindFailedLogs(ctx context.Context) ([]VCodeSmsLog, error) {
	var logs []VCodeSmsLog
	now := time.Now().Unix()
	err := s.db.WithContext(ctx).
		Where("status = ? AND CreateTime < ?", 0, now).
		Find(&logs).Error
	if err != nil {
		return nil, err
	} //如果status 为0 或者创建时间比现在要早,则发送错误信息
	return logs, nil
}

func (s *smsDao) Update(ctx context.Context, log VCodeSmsLog) error {
	log.UpdatedTime = time.Now().Unix() //更新时初始化时间戳
	return s.db.WithContext(ctx).Save(&log).Error
}
