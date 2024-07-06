package dao

import (
	"LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type CheckDAO interface {
	Create(ctx context.Context, check domain.Check) (int64, error)                         // 创建审核记录
	UpdateStatus(ctx context.Context, check domain.Check) error                            // 更新审核状态
	FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.CheckList, error) // 获取审核列表
	FindByID(ctx context.Context, checkId int64) (domain.Check, error)
	FindByPostId(ctx context.Context, postId int64) (domain.Check, error) // 获取审核详情
	GetCheckCount(ctx context.Context) (int64, error)
}

type checkDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

type Check struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`                     // 审核ID
	PostID    int64  `gorm:"not null"`                                     // 帖子ID
	Content   string `gorm:"type:text;not null"`                           // 审核内容
	Title     string `gorm:"size:255;not null"`                            // 审核标签
	Author    int64  `gorm:"column:author_id;index"`                       // 提交审核的用户ID
	Status    string `gorm:"size:20;not null;default:'Pending'"`           // 审核状态
	Remark    string `gorm:"type:text"`                                    // 审核备注
	CreatedAt int64  `gorm:"column:created_at;type:bigint;not null"`       // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint;not null;index"` // 更新时间
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

func (dao *checkDAO) FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.CheckList, error) {
	var checks []Check
	intSize := int(*pagination.Size)
	intOffset := int(*pagination.Offset)
	result := dao.db.WithContext(ctx).Limit(intSize).Offset(intOffset).Find(&checks)
	if result.Error != nil {
		dao.l.Error("failed to find all checks", zap.Error(result.Error))
		return nil, result.Error
	}
	// 将 Check 转换为 domain.Check
	var domainChecks []domain.CheckList
	for _, check := range checks {
		domainChecks = append(domainChecks, domain.CheckList{
			ID:        check.ID,
			PostID:    check.PostID,
			Title:     check.Title,
			UserID:    check.Author,
			Status:    check.Status,
			Remark:    check.Remark,
			UpdatedAt: check.UpdatedAt,
			CreatedAt: check.CreatedAt,
		})
	}
	return domainChecks, nil
}

func (dao *checkDAO) FindByID(ctx context.Context, checkId int64) (domain.Check, error) {
	var check Check
	result := dao.db.WithContext(ctx).Where("id = ?", checkId).First(&check)
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

func (dao *checkDAO) FindByPostId(ctx context.Context, postId int64) (domain.Check, error) {
	var check Check
	result := dao.db.WithContext(ctx).Where("post_id = ?", postId).First(&check)
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

func (dao *checkDAO) GetCheckCount(ctx context.Context) (int64, error) {
	var count int64
	if err := dao.db.WithContext(ctx).Model(&Check{}).Where("status = ?", constants.PostUnderReview).Count(&count).Error; err != nil {
		dao.l.Error("failed to get check count", zap.Error(err))
		return -1, err
	}
	return count, nil
}
