package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"context"
	"go.uber.org/zap"
)

type CheckRepository interface {
	Create(ctx context.Context, check domain.Check) (int64, error)                     // 创建审核记录
	UpdateStatus(ctx context.Context, check domain.Check) error                        // 更新审核状态
	FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) // 获取审核列表
	FindByID(ctx context.Context, checkID int64) (domain.Check, error)                 // 获取审核详情
}

type checkRepository struct {
	dao dao.CheckDAO
	l   *zap.Logger
}

func NewCheckRepository(dao dao.CheckDAO, l *zap.Logger) CheckRepository {
	return &checkRepository{
		dao: dao,
		l:   l,
	}
}

func (r *checkRepository) Create(ctx context.Context, check domain.Check) (int64, error) {
	return r.dao.Create(ctx, check)
}

func (r *checkRepository) UpdateStatus(ctx context.Context, check domain.Check) error {
	return r.dao.UpdateStatus(ctx, check)
}

func (r *checkRepository) FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) {
	return r.dao.FindAll(ctx, pagination)
}

func (r *checkRepository) FindByID(ctx context.Context, checkID int64) (domain.Check, error) {
	return r.dao.FindByID(ctx, checkID)
}
