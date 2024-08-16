package dao

import (
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

// VCodeSmsLog 表示用户认证操作的日志记录
type VCodeSmsLog struct {
	Id          int64  `gorm:"column:id;primaryKey;autoIncrement"`           // 自增ID
	SmsId       int64  `gorm:"column:sms_id"`                                // 短信类型ID
	SmsType     string `gorm:"column:sms_type"`                              // 短信类型
	Mobile      string `gorm:"column:mobile"`                                // 手机号
	VCode       string `gorm:"column:v_code"`                                // 验证码
	Driver      string `gorm:"column:driver"`                                // 服务商类型
	Status      int64  `gorm:"column:status"`                                // 发送状态，1为成功，0为失败
	StatusCode  string `gorm:"column:status_code"`                           // 状态码
	CreateTime  int64  `gorm:"column:created_at;type:bigint;not null"`       // 创建时间
	UpdatedTime int64  `gorm:"column:updated_at;type:bigint;not null;index"` // 更新时间
	DeletedTime int64  `gorm:"column:deleted_at;type:bigint;index"`          // 删除时间
}

func NewSmsDAO(db *gorm.DB, l *zap.Logger) SmsDAO {
	return &smsDao{
		db: db,
		l:  l,
	}
}

func (s *smsDao) Insert(ctx context.Context, log VCodeSmsLog) error {
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
