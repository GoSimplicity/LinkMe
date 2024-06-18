package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PlateRepository interface {
	CreatePlate(ctx *gin.Context, plate domain.Plate) error
	ListPlate(ctx *gin.Context, pagination domain.Pagination) ([]domain.Plate, error)
	UpdatePlate(ctx *gin.Context, plateId int64, uid int64) error
	DeletePlate(ctx *gin.Context, plateId int64, uid int64) error
}

type plateRepository struct {
	l   *zap.Logger
	dao dao.PlateDAO
}

func NewPlateRepository(l *zap.Logger, dao dao.PlateDAO) PlateRepository {
	return &plateRepository{
		l:   l,
		dao: dao,
	}
}

func (p *plateRepository) CreatePlate(ctx *gin.Context, plate domain.Plate) error {
	return p.dao.CreatePlate(ctx, plate)
}

func (p *plateRepository) ListPlate(ctx *gin.Context, pagination domain.Pagination) ([]domain.Plate, error) {
	//TODO implement me
	panic("implement me")
}

func (p *plateRepository) UpdatePlate(ctx *gin.Context, plateId int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (p *plateRepository) DeletePlate(ctx *gin.Context, plateId int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}
