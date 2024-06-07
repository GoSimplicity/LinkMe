package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
	"go.uber.org/zap"
)

type CheckService interface {
	SubmitCheck(ctx context.Context, check domain.Check) (int64, error)                   // 提交审核
	ApproveCheck(ctx context.Context, checkID int64, remark string) error                 // 审核通过
	RejectCheck(ctx context.Context, checkID int64, reason string) error                  // 审核拒绝
	ListChecks(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) // 获取审核列表
	CheckDetail(ctx context.Context, checkID int64) (domain.Check, error)                 // 获取审核详情
}

type checkService struct {
	repo repository.CheckRepository
	l    *zap.Logger
}

func NewCheckService(repo repository.CheckRepository, l *zap.Logger) CheckService {
	return &checkService{
		repo: repo,
		l:    l,
	}
}

func (s *checkService) SubmitCheck(ctx context.Context, check domain.Check) (int64, error) {
	panic("")
}

func (s *checkService) ApproveCheck(ctx context.Context, checkID int64, remark string) error {
	panic("")
}

func (s *checkService) RejectCheck(ctx context.Context, checkID int64, reason string) error {
	panic("")
}

func (s *checkService) ListChecks(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) {
	panic("")
}

func (s *checkService) CheckDetail(ctx context.Context, checkID int64) (domain.Check, error) {
	panic("")
}
