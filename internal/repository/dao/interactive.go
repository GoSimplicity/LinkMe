package dao

import (
	. "LinkMe/internal/repository/models"
	"context"
	"gorm.io/gorm"
)

// InteractiveDAO 互动数据访问对象接口
type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	BatchIncrReadCnt(ctx context.Context, biz []string, bizIds []int64) error
	InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error)
	GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error)
	Get(ctx context.Context, biz string, id int64) (Interactive, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error)
}

type interactiveDAO struct {
	db *gorm.DB
}

func NewInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &interactiveDAO{db: db}
}

func (i interactiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	//TODO implement me
	panic("implement me")
}

func (i interactiveDAO) BatchIncrReadCnt(ctx context.Context, biz []string, bizIds []int64) error {
	//TODO implement me
	panic("implement me")
}

func (i interactiveDAO) InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (i interactiveDAO) DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (i interactiveDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	//TODO implement me
	panic("implement me")
}

func (i interactiveDAO) GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (i interactiveDAO) GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (i interactiveDAO) Get(ctx context.Context, biz string, id int64) (Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (i interactiveDAO) GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error) {
	//TODO implement me
	panic("implement me")
}
