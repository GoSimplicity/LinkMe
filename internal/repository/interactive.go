package repository

import (
	. "LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/repository/models"
	"context"
	"errors"

	"go.uber.org/zap"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	// BatchIncrReadCnt biz 和 bizId 长度必须一致
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
	return c.dao.InsertLikeInfo(ctx, models.UserLikeBiz{
		BizName: biz,
		BizID:   id,
		Uid:     uid,
	})
}

func (c *CachedInteractiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	return c.dao.DeleteLikeInfo(ctx, models.UserLikeBiz{
		BizName: biz,
		BizID:   id,
		Uid:     uid,
	})
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
	ic, err := c.dao.Get(ctx, biz, id)
	if err != nil {
		c.l.Error(PostGetInteractiveERROR, zap.Error(err))
		return domain.Interactive{}, err
	}
	return toDomain(ic), nil
}

func (c *CachedInteractiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, id, uid)
	switch {
	case err == nil:
		// 如果没有错误，说明找到了点赞记录，返回true
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		// 如果错误是ErrRecordNotFound，说明没有找到点赞记录
		c.l.Error(PostGetLikedERROR, zap.Error(err))
		return false, nil
	default:
		// 如果是其他错误，返回false和错误信息
		return false, err
	}
}

func (c *CachedInteractiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectInfo(ctx, biz, id, uid)
	switch {
	case err == nil:
		// 如果没有错误，说明找到了收藏记录，返回true
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		// 如果错误是ErrRecordNotFound，说明没有找到收藏记录
		c.l.Error(PostGetCollectERROR, zap.Error(err))
		return false, nil
	default:
		// 如果是其他错误，返回false和错误信息
		return false, err
	}
}

func (c *CachedInteractiveRepository) GetById(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func toDomain(ic models.Interactive) domain.Interactive {
	return domain.Interactive{
		BizID:        ic.BizID,
		ReadCount:    ic.ReadCount,
		LikeCount:    ic.LikeCount,
		CollectCount: ic.CollectCount,
	}
}
