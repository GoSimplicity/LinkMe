package repository

import (
	"LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"context"
	"go.uber.org/zap"
	"log"
)

type CheckRepository interface {
	Create(ctx context.Context, check domain.Check) (int64, error)                         // 创建审核记录
	UpdateStatus(ctx context.Context, check domain.Check) error                            // 更新审核状态
	FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.CheckList, error) // 获取审核列表
	FindByID(ctx context.Context, checkID int64) (domain.Check, error)                     // 获取审核详情
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
	// 先查找是否存在该帖子审核信息
	dc, err := r.dao.FindByPostId(ctx, check.PostID)
	log.Println(dc)
	if dc.PostID != 0 && err == nil && dc.Status != constants.PostUnApproved {
		return 0, nil
	}
	// 创建新的审核信息
	id, err := r.dao.Create(ctx, check)
	if err != nil {
		r.l.Error("创建审核信息失败", zap.Error(err), zap.Error(err))
		return 0, err
	}
	return id, nil
}

func (r *checkRepository) UpdateStatus(ctx context.Context, check domain.Check) error {
	return r.dao.UpdateStatus(ctx, check)
}

func (r *checkRepository) FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.CheckList, error) {
	return r.dao.FindAll(ctx, pagination)
}

func (r *checkRepository) FindByID(ctx context.Context, checkID int64) (domain.Check, error) {
	return r.dao.FindByID(ctx, checkID)
}
