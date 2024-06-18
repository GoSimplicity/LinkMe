package dao

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PlateDAO interface {
	CreatePlate(ctx *gin.Context, plate domain.Plate) error
	ListPlate(ctx *gin.Context, pagination domain.Pagination) ([]domain.Plate, error)
	UpdatePlate(ctx *gin.Context, plateId int64, uid int64) error
	DeletePlate(ctx *gin.Context, plateId int64, uid int64) error
}

type plateDAO struct {
	l  *zap.Logger
	db *gorm.DB
}

func NewPlateDAO(l *zap.Logger, db *gorm.DB) PlateDAO {
	return &plateDAO{
		l:  l,
		db: db,
	}
}
func (p *plateDAO) CreatePlate(ctx *gin.Context, plate domain.Plate) error {
	newPlate := &models.Plate{
		Name:        plate.Name,
		Description: plate.Description,
		Uid:         plate.Uid,
	}
	err := p.db.WithContext(ctx).Create(newPlate).Error
	if err != nil {
		p.l.Error("Failed to create plate", zap.String("name", plate.Name), zap.Error(err))
		return err
	}
	return nil
}

func (p *plateDAO) ListPlate(ctx *gin.Context, pagination domain.Pagination) ([]domain.Plate, error) {
	//TODO implement me
	panic("implement me")
}

func (p *plateDAO) UpdatePlate(ctx *gin.Context, plateId int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (p *plateDAO) DeletePlate(ctx *gin.Context, plateId int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}
