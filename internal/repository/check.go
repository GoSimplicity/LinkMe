package repository

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

type CheckRepository interface {
	Create(ctx context.Context, check domain.Check) (int64, error)
	UpdateStatus(ctx context.Context, check domain.Check) error
	FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error)
	FindByID(ctx context.Context, checkID int64) (domain.Check, error)
	FindByPostId(ctx context.Context, postID uint) (domain.Check, error)
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

// Create 创建审核记录
func (r *checkRepository) Create(ctx context.Context, check domain.Check) (int64, error) {
	// 创建新的审核信息
	id, err := r.dao.Create(ctx, toDAOCheck(check))
	if err != nil {
		return -1, err
	}

	return id, nil
}

// UpdateStatus 更新审核状态
func (r *checkRepository) UpdateStatus(ctx context.Context, check domain.Check) error {
	return r.dao.UpdateStatus(ctx, toDAOCheck(check))
}

// FindAll 获取审核列表
func (r *checkRepository) FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) {
	checks, err := r.dao.FindAll(ctx, pagination)
	if err != nil {
		return nil, err
	}
	return toDomainChecks(checks), nil
}

// FindByID 获取审核详情
func (r *checkRepository) FindByID(ctx context.Context, checkID int64) (domain.Check, error) {
	check, err := r.dao.FindByID(ctx, checkID)
	if err != nil {
		return domain.Check{}, err
	}
	return toDomainCheck(check), nil
}

// FindByPostId 根据帖子ID获取审核信息
func (r *checkRepository) FindByPostId(ctx context.Context, postID uint) (domain.Check, error) {
	check, err := r.dao.FindByPostId(ctx, postID)
	if err != nil {
		return domain.Check{}, err
	}
	return toDomainCheck(check), nil
}

// toDAOCheck 将 domain.Check 转换为 dao.Check
func toDAOCheck(domainCheck domain.Check) dao.Check {
	return dao.Check{
		ID:        domainCheck.ID,
		PostID:    domainCheck.PostID,
		Content:   domainCheck.Content,
		Title:     domainCheck.Title,
		PlateID:   domainCheck.PlateID,
		Uid:       domainCheck.Uid,
		Status:    domainCheck.Status,
		Remark:    domainCheck.Remark,
		CreatedAt: domainCheck.CreatedAt,
		UpdatedAt: domainCheck.UpdatedAt,
	}
}

// toDomainCheck 将 dao.Check 转换为 domain.Check
func toDomainCheck(daoCheck dao.Check) domain.Check {
	return domain.Check{
		ID:        daoCheck.ID,
		PostID:    daoCheck.PostID,
		Content:   daoCheck.Content,
		Title:     daoCheck.Title,
		Uid:       daoCheck.Uid,
		PlateID:   daoCheck.PlateID,
		Status:    daoCheck.Status,
		Remark:    daoCheck.Remark,
		CreatedAt: daoCheck.CreatedAt,
		UpdatedAt: daoCheck.UpdatedAt,
	}
}

// toDomainChecks 将 []dao.Check 转换为 []domain.Check
func toDomainChecks(daoChecks []dao.Check) []domain.Check {
	domainChecks := make([]domain.Check, len(daoChecks))
	for i, daoCheck := range daoChecks {
		domainChecks[i] = toDomainCheck(daoCheck)
	}
	return domainChecks
}
