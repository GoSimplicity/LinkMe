package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/repository/models"
	"context"
	"github.com/ecodeclub/ekit/sqlx"
)

// ErrWaitingSMSNotFound 错误常量定义
var ErrWaitingSMSNotFound = dao.ErrWaitingSMSNotFound

// SendCodeRepository 接口定义了异步 SMS 记录操作的相关方法
type SendCodeRepository interface {
	Add(ctx context.Context, s domain.AsyncSms) error
	PreemptWaitingSMS(ctx context.Context) (domain.AsyncSms, error)
	ReportScheduleResult(ctx context.Context, id int64, success bool) error
}

// sendCodeRepository 实现了 AsyncSmsRepository 接口
type sendCodeRepository struct {
	dao dao.SendCodeDAO
}

// NewAsyncSMSRepository 创建并返回一个新的 asyncSmsRepository 实例
func NewAsyncSMSRepository(dao dao.SendCodeDAO) SendCodeRepository {
	return &sendCodeRepository{
		dao: dao,
	}
}

// Add 添加一个异步 SMS 记录
func (a *sendCodeRepository) Add(ctx context.Context, s domain.AsyncSms) error {
	return a.dao.Insert(ctx, models.Sms{
		Config: sqlx.JsonColumn[models.SmsConfig]{
			Val: models.SmsConfig{
				TplId:   s.TplId,
				Args:    s.Args,
				Numbers: s.Numbers,
			},
			Valid: true,
		},
		RetryMax: s.RetryMax,
	})
}

// PreemptWaitingSMS 抢占一个待发送的 SMS 记录
func (a *sendCodeRepository) PreemptWaitingSMS(ctx context.Context) (domain.AsyncSms, error) {
	as, err := a.dao.GetWaitingSMS(ctx)
	if err != nil {
		return domain.AsyncSms{}, err
	}
	return domain.AsyncSms{
		Id:       as.Id,
		TplId:    as.Config.Val.TplId,
		Numbers:  as.Config.Val.Numbers,
		Args:     as.Config.Val.Args,
		RetryMax: as.RetryMax,
	}, nil
}

// ReportScheduleResult 报告调度结果
func (a *sendCodeRepository) ReportScheduleResult(ctx context.Context, id int64, success bool) error {
	if success {
		return a.dao.MarkSuccess(ctx, id)
	}
	return a.dao.MarkFailed(ctx, id)
}
