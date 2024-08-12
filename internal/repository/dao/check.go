package dao

import (
	"context"
	"errors"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type CheckDAO interface {
	Create(ctx context.Context, check Check) (int64, error)                     // 创建审核记录
	UpdateStatus(ctx context.Context, check Check) error                        // 更新审核状态
	FindAll(ctx context.Context, pagination domain.Pagination) ([]Check, error) // 获取审核列表
	FindByID(ctx context.Context, checkId int64) (Check, error)
	FindByPostId(ctx context.Context, postId uint) (Check, error) // 获取审核详情
	GetCheckCount(ctx context.Context) (int64, error)
}

type checkDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

type Check struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`                     // 审核ID
	PostID    uint   `gorm:"not null"`                                     // 帖子ID
	Content   string `gorm:"type:text;not null"`                           // 审核内容
	Title     string `gorm:"size:255;not null"`                            // 审核标签
	Author    int64  `gorm:"column:author_id;index"`                       // 提交审核的用户ID
	Status    uint8  `gorm:"default:0"`                                    // 审核状态
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

func (dao *checkDAO) Create(ctx context.Context, check Check) (int64, error) {
	now := time.Now().UnixMilli()

	// 判断传入的check是否有效
	if check.PostID == 0 || check.Content == "" || check.Title == "" || check.Author == 0 {
		return 0, errors.New("invalid input: missing required fields")
	}

	check.CreatedAt = now
	check.UpdatedAt = now

	if err := dao.db.WithContext(ctx).Create(&check).Error; err != nil {
		dao.l.Error("failed to create check", zap.Error(err))
		return 0, err
	}

	return check.ID, nil
}

func (dao *checkDAO) UpdateStatus(ctx context.Context, check Check) error {
	if check.ID == 0 {
		return errors.New("invalid input: missing required fields")
	}

	result := dao.db.WithContext(ctx).Model(&Check{}).Where("id = ?", check.ID).Updates(Check{
		Status:    check.Status,
		UpdatedAt: time.Now().UnixMilli(),
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

func (dao *checkDAO) FindAll(ctx context.Context, pagination domain.Pagination) ([]Check, error) {
	var checks []Check

	intSize := int(*pagination.Size)
	intOffset := int(*pagination.Offset)

	result := dao.db.WithContext(ctx).
		Limit(intSize).
		Offset(intOffset).
		Find(&checks)

	if result.Error != nil {
		dao.l.Error("failed to find all checks", zap.Error(result.Error))
		return nil, result.Error
	}

	return checks, nil
}

func (dao *checkDAO) FindByID(ctx context.Context, checkId int64) (Check, error) {
	var check Check

	result := dao.db.WithContext(ctx).Where("id = ?", checkId).First(&check)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Check{}, nil
		}
		dao.l.Error("failed to find check by ID", zap.Error(result.Error))
		return Check{}, result.Error
	}

	return check, nil
}

func (dao *checkDAO) FindByPostId(ctx context.Context, postId uint) (Check, error) {
	var check Check

	result := dao.db.WithContext(ctx).Where("post_id = ?", postId).First(&check)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Check{}, nil
		}
		dao.l.Error("failed to find check by post ID", zap.Error(result.Error))
		return Check{}, result.Error
	}

	return check, nil
}

func (dao *checkDAO) GetCheckCount(ctx context.Context) (int64, error) {
	var count int64

	if err := dao.db.WithContext(ctx).Model(&Check{}).Where("status = ?", domain.UnderReview).Count(&count).Error; err != nil {
		dao.l.Error("failed to get check count", zap.Error(err))
		return -1, err
	}

	return count, nil
}
