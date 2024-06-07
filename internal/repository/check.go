package repository

import (
	"LinkMe/internal/domain"
	"context"
	"database/sql"
	"go.uber.org/zap"
)

type CheckRepository interface {
	Create(ctx context.Context, check domain.Check) (int64, error)                     // 创建审核记录
	UpdateStatus(ctx context.Context, check domain.Check) error                        // 更新审核状态
	FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) // 获取审核列表
	FindByID(ctx context.Context, checkID int64) (domain.Check, error)                 // 获取审核详情
}

type CheckRepositoryImpl struct {
	db *sql.DB
	l  *zap.Logger
}

func NewCheckRepository(db *sql.DB, l *zap.Logger) CheckRepository {
	return &CheckRepositoryImpl{
		db: db,
		l:  l,
	}
}

func (r *CheckRepositoryImpl) Create(ctx context.Context, check domain.Check) (int64, error) {
	panic("")
}

func (r *CheckRepositoryImpl) UpdateStatus(ctx context.Context, check domain.Check) error {
	panic("")
}

func (r *CheckRepositoryImpl) FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) {
	panic("")
}

func (r *CheckRepositoryImpl) FindByID(ctx context.Context, checkID int64) (domain.Check, error) {
	panic("")
}
