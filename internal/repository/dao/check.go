package dao

import (
	"LinkMe/internal/domain"
	. "LinkMe/internal/repository/models"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type CheckDAO interface {
	Create(ctx context.Context, check domain.Check) (int64, error)                     // 创建审核记录
	UpdateStatus(ctx context.Context, check domain.Check) error                        // 更新审核状态
	FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) // 获取审核列表
	FindByID(ctx context.Context, checkID int64) (domain.Check, error)                 // 获取审核详情
}

type checkDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

func NewCheckDAO(db *gorm.DB, l *zap.Logger) CheckDAO {
	return &checkDAO{
		db: db,
		l:  l,
	}
}

func (dao *checkDAO) Create(ctx context.Context, check domain.Check) (int64, error) {
	// 判断传入的check是否有效
	if check.PostID == 0 || check.Content == "" || check.Title == "" || check.UserID == 0 {
		return 0, errors.New("invalid input: missing required fields")
	}
	now := time.Now().UnixMilli()
	modelCheck := Check{
		PostID:    check.PostID,
		Content:   check.Content,
		Title:     check.Title,
		Author:    check.UserID,
		Status:    check.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	// 开启事务
	tx := dao.db.WithContext(ctx).Begin()
	if err := tx.Create(&modelCheck).Error; err != nil {
		// 如果出现错误直接回滚
		tx.Rollback()
		dao.l.Error("failed to create check", zap.Error(err))
		return 0, err
	}
	// 提交事务
	tx.Commit()
	return modelCheck.ID, nil
}

func (dao *checkDAO) UpdateStatus(ctx context.Context, check domain.Check) error {
	// 判断传入的check是否有效
	if check.ID == 0 || check.Status == "" {
		return errors.New("invalid input: missing required fields")
	}
	now := time.Now().UnixMilli()
	// 更新状态和更新时间
	result := dao.db.WithContext(ctx).Model(&Check{}).Where("id = ?", check.ID).Updates(Check{
		Status:    check.Status,
		UpdatedAt: now,
	})
	if result.Error != nil {
		dao.l.Error("failed to update check status", zap.Error(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records updated")
	}
	return nil
}

func (dao *checkDAO) FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) {
	var checks []Check
	intSize := int(*pagination.Size)
	intOffset := int(*pagination.Offset)
	result := dao.db.WithContext(ctx).Limit(intSize).Offset(intOffset).Find(&checks)
	if result.Error != nil {
		dao.l.Error("failed to find all checks", zap.Error(result.Error))
		return nil, result.Error
	}
	// 将 models.Check 转换为 domain.Check
	var domainChecks []domain.Check
	for _, check := range checks {
		domainChecks = append(domainChecks, domain.Check{
			ID:      check.ID,
			PostID:  check.PostID,
			Content: check.Content,
			Title:   check.Title,
			UserID:  check.Author,
			Status:  check.Status,
			Remark:  check.Remark,
		})
	}
	return domainChecks, nil
}

func (dao *checkDAO) FindByID(ctx context.Context, checkID int64) (domain.Check, error) {
	var check Check
	result := dao.db.WithContext(ctx).First(&check, checkID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.Check{}, nil
		}
		dao.l.Error("failed to find check by ID", zap.Error(result.Error))
		return domain.Check{}, result.Error
	}
	return domain.Check{
		ID:      check.ID,
		PostID:  check.PostID,
		Content: check.Content,
		Title:   check.Title,
		UserID:  check.Author,
		Status:  check.Status,
		Remark:  check.Remark,
	}, nil
}
