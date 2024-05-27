package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/repository/models"
	"context"

	"go.uber.org/zap"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	// biz 和 bizId 长度必须一致
	BatchIncrReadCnt(ctx context.Context, biz []string, id []int64) error                     // 批量更新阅读计数(与kafka配合使用)
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error                      // 增加阅读计数
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error                      // 减少阅读计数
	IncrCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error // 收藏
	DecrCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error // 取消收藏
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)                // 获取互动信息
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)                 // 检查是否已点赞
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)             // 检查是否被收藏
	GetById(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error)       // 批量获取互动信息
}

type CachedInteractiveRepository struct {
	dao dao.InteractiveDAO
	l   *zap.Logger
}

func NewInteractiveRepository(dao dao.InteractiveDAO, l *zap.Logger) InteractiveRepository {
	return &CachedInteractiveRepository{dao: dao, l: l}
}

func (c *CachedInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, id int64) error {
	return c.dao.IncrReadCnt(ctx, biz, id)
}

func (c *CachedInteractiveRepository) BatchIncrReadCnt(ctx context.Context, biz []string, id []int64) error {
	//TODO implement me
	panic("implement me")
}

func (c *CachedInteractiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (c *CachedInteractiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (c *CachedInteractiveRepository) IncrCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	return c.dao.InsertCollectionBiz(ctx, models.UserCollectionBiz{
		BizName:      biz,
		BizID:        id,
		CollectionId: cid,
		Uid:          uid,
	})
}
func (c *CachedInteractiveRepository) DecrCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	return c.dao.DeleteCollectionBiz(ctx, models.UserCollectionBiz{
		BizName:      biz,
		BizID:        id,
		CollectionId: cid,
		Uid:          uid,
	})
}

func (c *CachedInteractiveRepository) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CachedInteractiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CachedInteractiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CachedInteractiveRepository) GetById(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}
