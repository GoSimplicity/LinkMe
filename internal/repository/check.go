package repository

import (
	"context"
	"fmt"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

type CheckRepository interface {
	Create(ctx context.Context, check domain.Check) (int64, error)                     // 创建审核记录
	BatchCreate(ctx context.Context, checks []domain.Check) error                      // 批量创建审核记录
	UpdateStatus(ctx context.Context, check domain.Check) error                        // 更新审核状态
	FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) // 获取审核列表
	FindByID(ctx context.Context, checkID int64) (domain.Check, error)                 // 获取审核详情
	FindByPostId(ctx context.Context, postID uint) (domain.Check, error)               // 根据帖子ID获取审核信息
	GetCheckCount(ctx context.Context) (int64, error)                                  // 获取审核数量
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

func (r *checkRepository) FindByPostId(ctx context.Context, postID uint) (domain.Check, error) {
	check, err := r.dao.FindByPostId(ctx, postID)
	if err != nil {
		return domain.Check{}, err
	}
	return toDomainCheck(check), nil
}

func (r *checkRepository) Create(ctx context.Context, check domain.Check) (int64, error) {
	// 先查找是否存在该帖子审核信息
	dc, err := r.dao.FindByPostId(ctx, check.PostID)
	if dc.PostID != 0 && err == nil {
		return -1, nil
	}

	// 创建新的审核信息
	id, err := r.dao.Create(ctx, toDAOCheck(check))
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r *checkRepository) BatchCreate(ctx context.Context, checks []domain.Check) error {
	// 将 domain.Check 转换为 dao.Check
	daoChecks := make([]dao.Check, len(checks))
	for i, check := range checks {
		daoChecks[i] = toDAOCheck(check)
	}

	// 防止重复创建
	for _, check := range checks {
		if check.PostID == 0 {
			return fmt.Errorf("post_id is required")
		}

		dc, err := r.dao.FindByPostId(ctx, check.PostID)
		if err != nil {
			return err
		}
		if dc.PostID != 0 {
			return nil
		}
	}

	// 批量创建审核记录
	if err := r.dao.BatchCreate(ctx, daoChecks); err != nil {
		return err
	}

	return nil
}

func (r *checkRepository) UpdateStatus(ctx context.Context, check domain.Check) error {
	return r.dao.UpdateStatus(ctx, toDAOCheck(check))
}

func (r *checkRepository) FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) {
	checks, err := r.dao.FindAll(ctx, pagination)
	if err != nil {
		return nil, err
	}
	return toDomainChecks(checks), nil
}

func (r *checkRepository) FindByID(ctx context.Context, checkID int64) (domain.Check, error) {
	check, err := r.dao.FindByID(ctx, checkID)
	if err != nil {
		return domain.Check{}, err
	}
	return toDomainCheck(check), nil
}

func (r *checkRepository) GetCheckCount(ctx context.Context) (int64, error) {
	return r.dao.GetCheckCount(ctx)
}

// toDAOCheck 将 domain.Check 转换为 dao.Check
func toDAOCheck(domainCheck domain.Check) dao.Check {
	return dao.Check{
		ID:        domainCheck.ID,
		PostID:    domainCheck.PostID,
		Content:   domainCheck.Content,
		Title:     domainCheck.Title,
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
