package dao

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type PlateDAO interface {
	CreatePlate(ctx context.Context, plate domain.Plate) error
	ListPlate(ctx context.Context, pagination domain.Pagination) ([]Plate, error)
	UpdatePlate(ctx context.Context, plate domain.Plate) error
	DeletePlate(ctx context.Context, plateId int64, uid int64) error
}

type plateDAO struct {
	l  *zap.Logger
	db *gorm.DB
}

type Plate struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`      // 板块ID
	Name        string `gorm:"size:255;not null;uniqueIndex"` // 板块名称
	Description string `gorm:"type:text"`                     // 板块描述
	CreateTime  int64  `gorm:"column:created_at;type:bigint"` // 创建时间
	UpdatedTime int64  `gorm:"column:updated_at;type:bigint"` // 更新时间
	DeletedTime int64  `gorm:"column:deleted_at;type:bigint"` // 删除时间
	Deleted     bool   `gorm:"column:deleted;default:false"`  // 是否删除
	Uid         int64  `gorm:"index"`                         // 板主id
	Posts       []Post `gorm:"foreignKey:PlateID"`            // 帖子关系
}

func NewPlateDAO(l *zap.Logger, db *gorm.DB) PlateDAO {
	return &plateDAO{
		l:  l,
		db: db,
	}
}

func (p *plateDAO) CreatePlate(ctx context.Context, plate domain.Plate) error {
	now := time.Now().UnixMilli()
	newPlate := &Plate{
		Name:        plate.Name,
		Description: plate.Description,
		Uid:         plate.Uid,
		CreateTime:  now,
		UpdatedTime: now,
	}

	if er := p.db.WithContext(ctx).Create(newPlate).Error; er != nil {
		p.l.Error("Failed to create plate", zap.String("name", plate.Name), zap.Error(er))
		return er
	}

	return nil
}

func (p *plateDAO) ListPlate(ctx context.Context, pagination domain.Pagination) ([]Plate, error) {
	var plates []Plate

	intSize := int(*pagination.Size)
	intOffset := int(*pagination.Offset)
	err := p.db.WithContext(ctx).
		Where("deleted = ?", false).
		Limit(intSize).
		Offset(intOffset).
		Find(&plates).Error
	if err != nil {
		p.l.Error("Failed to list plates", zap.Error(err))
		return nil, err
	}

	return plates, nil
}

func (p *plateDAO) UpdatePlate(ctx context.Context, plate domain.Plate) error {
	now := time.Now().UnixMilli()
	updateData := map[string]interface{}{
		"name":        plate.Name,
		"description": plate.Description,
		"updated_at":  now,
	}

	if err := p.db.WithContext(ctx).Model(&Plate{}).
		Where("id = ? AND uid = ?", plate.ID, plate.Uid).
		Updates(updateData).Error; err != nil {
		p.l.Error("Failed to update plate", zap.Error(err))
		return err
	}

	return nil
}

func (p *plateDAO) DeletePlate(ctx context.Context, plateId int64, uid int64) error {
	// 进行软删除
	now := time.Now().UnixMilli()
	updateData := map[string]interface{}{
		"deleted":    true,
		"deleted_at": now,
		"updated_at": now,
	}

	if err := p.db.WithContext(ctx).Model(&Plate{}).
		Where("id = ? AND uid = ?", plateId, uid).
		Updates(updateData).Error; err != nil {
		p.l.Error("Failed to soft delete plate", zap.Error(err))
		return err
	}

	return nil
}
