package repository

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

type PlateRepository interface {
	CreatePlate(ctx context.Context, plate domain.Plate) error
	ListPlate(ctx context.Context, pagination domain.Pagination) ([]domain.Plate, error)
	UpdatePlate(ctx context.Context, plate domain.Plate) error
	DeletePlate(ctx context.Context, plateId int64, uid int64) error
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

func (p *plateRepository) CreatePlate(ctx context.Context, plate domain.Plate) error {
	return p.dao.CreatePlate(ctx, plate)
}

func (p *plateRepository) ListPlate(ctx context.Context, pagination domain.Pagination) ([]domain.Plate, error) {
	plates, err := p.dao.ListPlate(ctx, pagination)
	if err != nil {
		return nil, err
	}
	return fromDomainSlicePlate(plates), err
}

func (p *plateRepository) UpdatePlate(ctx context.Context, plate domain.Plate) error {
	return p.dao.UpdatePlate(ctx, plate)
}

func (p *plateRepository) DeletePlate(ctx context.Context, plateId int64, uid int64) error {
	return p.dao.DeletePlate(ctx, plateId, uid)
}

// 将dao层对象转为领域层对象
func fromDomainSlicePlate(post []dao.Plate) []domain.Plate {
	domainPlate := make([]domain.Plate, len(post))
	for i, repoPlate := range post {
		domainPlate[i] = domain.Plate{
			ID:          repoPlate.ID,
			Name:        repoPlate.Name,
			Uid:         repoPlate.Uid,
			Description: repoPlate.Description,
			CreatedAt:   repoPlate.CreateTime,
			UpdatedAt:   repoPlate.UpdatedTime,
			DeletedAt:   repoPlate.DeletedTime,
			Deleted:     repoPlate.Deleted,
		}
	}
	return domainPlate
}
