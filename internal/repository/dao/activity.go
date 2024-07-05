package dao

import (
	"LinkMe/internal/repository/models"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ActivityDAO interface {
	GetRecentActivity(ctx context.Context) (models.RecentActivity, error)
	SetRecentActivity(ctx context.Context, mr models.RecentActivity) error
}

type activityDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

func NewActivityDAO(db *gorm.DB, l *zap.Logger) ActivityDAO {
	return &activityDAO{
		db: db,
		l:  l,
	}
}

func (a *activityDAO) GetRecentActivity(ctx context.Context) (models.RecentActivity, error) {
	var mr models.RecentActivity
	if err := a.db.WithContext(ctx).Model(models.RecentActivity{}).Find(&mr).Error; err != nil {
		a.l.Error("get recent activity failed", zap.Error(err))
		return models.RecentActivity{}, err
	}
	return mr, nil
}

func (a *activityDAO) SetRecentActivity(ctx context.Context, mr models.RecentActivity) error {
	if err := a.db.WithContext(ctx).Create(&mr).Error; err != nil {
		a.l.Error("set recent activity failed", zap.Error(err))
		return err
	}
	return nil
}
