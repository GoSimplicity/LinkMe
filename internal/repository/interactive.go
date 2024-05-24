package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"context"
	"go.uber.org/zap"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	BatchIncrReadCnt(ctx context.Context, biz []string, bizId []int64) error // biz 和 bizId 长度必须一致
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	GetById(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error)
}
type CachedInteractiveRepository struct {
	dao dao.InteractiveDAO
	l   *zap.Logger
}

func NewInteractiveRepository(dao dao.InteractiveDAO, l *zap.Logger) InteractiveRepository {
	return &CachedInteractiveRepository{dao: dao, l: l}
}

func (c CachedInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	//TODO implement me
	panic("implement me")
}

func (c CachedInteractiveRepository) BatchIncrReadCnt(ctx context.Context, biz []string, bizId []int64) error {
	//TODO implement me
	panic("implement me")
}

func (c CachedInteractiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (c CachedInteractiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (c CachedInteractiveRepository) AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (c CachedInteractiveRepository) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (c CachedInteractiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c CachedInteractiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c CachedInteractiveRepository) GetById(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}
