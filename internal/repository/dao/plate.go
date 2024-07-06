package dao

import (
	"LinkMe/internal/domain"
	"context"
	"errors"
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
	var existingPlate Plate
	err := p.db.WithContext(ctx).Where("name = ? AND uid = ?", plate.Name, plate.Uid).First(&existingPlate).Error
	if err == nil {
		p.l.Warn("Plate already exists", zap.String("name", plate.Name), zap.Int64("uid", plate.Uid))
		return errors.New("plate already")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		p.l.Error("Failed to check if plate exists", zap.String("name", plate.Name), zap.Error(err))
		return err
	}
	now := time.Now().UnixMilli()
	newPlate := &Plate{
		Name:        plate.Name,
		Description: plate.Description,
		Uid:         plate.Uid,
		CreateTime:  now,
		UpdatedTime: now,
	}
	if er := p.db.WithContext(ctx).Create(newPlate).Error; er != nil {
		p.l.Error("Failed to create plate", zap.String("name", plate.Name), zap.Error(err))
		return er
	}
	return nil
}

func (p *plateDAO) ListPlate(ctx context.Context, pagination domain.Pagination) ([]Plate, error) {
	var plates []Plate
	intSize := int(*pagination.Size)
	intOffset := int(*pagination.Offset)
	err := p.db.WithContext(ctx).
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
	// 查询当前的 plate
	var existingPlate Plate
	if err := p.db.WithContext(ctx).Where("id = ? AND uid = ?", plate.ID, plate.Uid).First(&existingPlate).Error; err != nil {
		p.l.Error("Failed to find plate", zap.Int64("id", plate.ID), zap.Error(err))
		return err
	}
	// 检查是否有变化
	if existingPlate.Name == plate.Name && existingPlate.Description == plate.Description {
		p.l.Info("No changes detected, update skipped", zap.Int64("id", plate.ID))
		return errors.New("no changes detected")
	}
	// 更新数据
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
	var existingPlate Plate
	if err := p.db.WithContext(ctx).Where("id = ? AND uid = ?", plateId, uid).First(&existingPlate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.l.Warn("Plate not found", zap.Int64("id", plateId), zap.Int64("uid", uid))
			return errors.New("plate not found")
		}
		p.l.Error("Failed to find plate", zap.Int64("id", plateId), zap.Error(err))
		return err
	}

	// 检查是否已被删除
	if existingPlate.Deleted {
		p.l.Info("Plate already deleted", zap.Int64("id", plateId))
		return errors.New("plate already deleted")
	}
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
