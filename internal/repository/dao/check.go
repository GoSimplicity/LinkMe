package dao

import (
	"context"
	"errors"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CheckDAO interface {
	Create(ctx context.Context, check Check) (int64, error)
	UpdateStatus(ctx context.Context, check Check) error
	FindAll(ctx context.Context, pagination domain.Pagination) ([]Check, error)
	FindByID(ctx context.Context, checkId int64) (Check, error)
	FindByPostId(ctx context.Context, postId uint) (Check, error)
}

type checkDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

type Check struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"` // 审核ID
	PostID    uint   `gorm:"not null"`                 // 帖子ID
	Content   string `gorm:"type:text;not null"`       // 审核内容
	Title     string `gorm:"size:255;not null"`        // 审核标签
	BizId     int64  `gorm:"index:idx_biz_type_id"`    // 业务ID: Note:为了让审核模块复用(即既能审核帖子又能审核评论)，其中1：表示帖子业务，2：表示评论业务
	PlateID   int64  `gorm:"index"`
	Uid       int64  `gorm:"column:uid;index"`                             // 提交审核的用户ID
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

// Create 创建审核记录
func (dao *checkDAO) Create(ctx context.Context, check Check) (int64, error) {
	now := time.Now().UnixMilli()

	// 判断传入的check是否有效
	if check.PostID == 0 || check.Content == "" || (check.Title == "" && check.BizId == 1) || check.Uid == 0 {
		return 0, errors.New("无效输入：缺少必填字段")
	}

	check.CreatedAt = now
	check.UpdatedAt = now

	if err := dao.db.WithContext(ctx).Create(&check).Error; err != nil {
		dao.l.Error("创建审核记录失败", zap.Error(err))
		return 0, err
	}

	return check.ID, nil
}

// UpdateStatus 更新审核状态
func (dao *checkDAO) UpdateStatus(ctx context.Context, check Check) error {
	if check.ID == 0 {
		return errors.New("无效输入：缺少必填字段")
	}

	result := dao.db.WithContext(ctx).Model(&Check{}).Where("id = ?", check.ID).Updates(Check{
		Status:    check.Status,
		UpdatedAt: time.Now().UnixMilli(),
	})

	if result.Error != nil {
		dao.l.Error("更新审核状态失败", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("未更新任何记录")
	}

	return nil
}

// FindAll 获取审核列表
func (dao *checkDAO) FindAll(ctx context.Context, pagination domain.Pagination) ([]Check, error) {
	var checks []Check

	intSize := int(*pagination.Size)
	intOffset := int(*pagination.Offset)

	result := dao.db.WithContext(ctx).
		Limit(intSize).
		Offset(intOffset).
		Find(&checks)

	if result.Error != nil {
		dao.l.Error("获取所有审核记录失败", zap.Error(result.Error))
		return nil, result.Error
	}

	return checks, nil
}

// FindByID 获取审核详情
func (dao *checkDAO) FindByID(ctx context.Context, checkId int64) (Check, error) {
	var check Check

	result := dao.db.WithContext(ctx).Where("id = ?", checkId).First(&check)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Check{}, nil
		}
		dao.l.Error("根据ID查找审核记录失败", zap.Error(result.Error))
		return Check{}, result.Error
	}

	return check, nil
}

// FindByPostId 根据帖子ID获取审核信息
func (dao *checkDAO) FindByPostId(ctx context.Context, postId uint) (Check, error) {
	var check Check

	result := dao.db.WithContext(ctx).Where("post_id = ?", postId).First(&check)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Check{}, nil
		}
		dao.l.Error("根据帖子ID查找审核记录失败", zap.Error(result.Error))
		return Check{}, result.Error
	}

	return check, nil
}
