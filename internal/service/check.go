package service

import (
	"LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
	"go.uber.org/zap"
)

type CheckService interface {
	SubmitCheck(ctx context.Context, check domain.Check) (int64, error)                   // 提交审核
	ApproveCheck(ctx context.Context, checkID int64, remark string) error                 // 审核通过
	RejectCheck(ctx context.Context, checkID int64, remark string) error                  // 审核拒绝
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
	// 设置状态为审核中
	check.Status = constants.PostUnderReview
	s.l.Info("Submitting check", zap.Int64("PostID", check.PostID), zap.Int64("UserID", check.UserID))
	id, err := s.repo.Create(ctx, check)
	if err != nil {
		s.l.Error("Failed to submit check", zap.Error(err))
		return 0, err
	}
	return id, nil
}

func (s *checkService) ApproveCheck(ctx context.Context, checkID int64, remark string) error {
	// 更新审核状态为通过
	err := s.repo.UpdateStatus(ctx, domain.Check{
		ID:     checkID,
		Remark: remark,
		Status: constants.PostApproved,
	})
	if err != nil {
		s.l.Error("Failed to approve check", zap.Error(err))
		return err
	}
	s.l.Info("Approved check", zap.Int64("CheckID", checkID), zap.String("Remark", remark))
	return nil
}

func (s *checkService) RejectCheck(ctx context.Context, checkID int64, remark string) error {
	// 更新审核状态为拒绝
	err := s.repo.UpdateStatus(ctx, domain.Check{
		ID:     checkID,
		Remark: remark,
		Status: constants.PostUnApproved,
	})
	if err != nil {
		s.l.Error("Failed to reject check", zap.Error(err))
		return err
	}
	s.l.Info("Rejected check", zap.Int64("CheckID", checkID), zap.String("Remark", remark))
	return nil
}

func (s *checkService) ListChecks(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) {
	checks, err := s.repo.FindAll(ctx, pagination)
	if err != nil {
		s.l.Error("Failed to list checks", zap.Error(err))
		return nil, err
	}
	s.l.Info("Listed checks", zap.Int("Count", len(checks)))
	return checks, nil
}

func (s *checkService) CheckDetail(ctx context.Context, checkID int64) (domain.Check, error) {
	check, err := s.repo.FindByID(ctx, checkID)
	if err != nil {
		s.l.Error("Failed to get check detail", zap.Error(err))
		return domain.Check{}, err
	}
	s.l.Info("Fetched check detail", zap.Int64("CheckID", checkID))
	return check, nil
}
