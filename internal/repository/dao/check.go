package dao

import (
	"LinkMe/internal/domain"
	"context"
	"database/sql"
	"go.uber.org/zap"
)

type CheckDAO interface {
	Create(ctx context.Context, check domain.Check) (int64, error)                     // 创建审核记录
	UpdateStatus(ctx context.Context, check domain.Check) error                        // 更新审核状态
	FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) // 获取审核列表
	FindByID(ctx context.Context, checkID int64) (domain.Check, error)                 // 获取审核详情
}

type CheckDAOImpl struct {
	db *sql.DB
	l  *zap.Logger
}

func NewCheckDAO(db *sql.DB, l *zap.Logger) CheckDAO {
	return &CheckDAOImpl{
		db: db,
		l:  l,
	}
}

func (dao *CheckDAOImpl) Create(ctx context.Context, check domain.Check) (int64, error) {
	panic("")
}

func (dao *CheckDAOImpl) UpdateStatus(ctx context.Context, check domain.Check) error {
	panic("")
}

func (dao *CheckDAOImpl) FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) {
	panic("")
}

func (dao *CheckDAOImpl) FindByID(ctx context.Context, checkID int64) (domain.Check, error) {
	panic("")
}
