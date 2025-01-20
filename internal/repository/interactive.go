package repository

import (
	"context"
	"errors"

	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"

	"go.uber.org/zap"
)

type InteractiveRepository interface {
	BatchIncrReadCnt(ctx context.Context, postIds []uint) error
	IncrReadCnt(ctx context.Context, postId uint) error
	IncrLike(ctx context.Context, postId uint, uid int64) error
	DecrLike(ctx context.Context, postId uint, uid int64) error
	IncrCollectionItem(ctx context.Context, postId uint, uid int64) error
	DecrCollectionItem(ctx context.Context, postId uint, uid int64) error
	Get(ctx context.Context, postId uint) (domain.Interactive, error)
	Liked(ctx context.Context, postId uint, uid int64) (bool, error)
	Collected(ctx context.Context, postId uint, uid int64) (bool, error)
	GetById(ctx context.Context, postIds []uint) ([]domain.Interactive, error)
}

type InteractiveRepositoryImpl struct {
	dao dao.InteractiveDAO
	l   *zap.Logger
}

func NewInteractiveRepository(dao dao.InteractiveDAO, l *zap.Logger) InteractiveRepository {
	return &InteractiveRepositoryImpl{
		dao: dao,
		l:   l,
	}
}

// BatchIncrReadCnt 批量增加阅读计数
func (i *InteractiveRepositoryImpl) BatchIncrReadCnt(ctx context.Context, postIds []uint) error {
	return i.dao.BatchIncrReadCnt(ctx, postIds)
}

// IncrReadCnt 增加阅读计数
func (i *InteractiveRepositoryImpl) IncrReadCnt(ctx context.Context, postId uint) error {
	return i.dao.IncrReadCnt(ctx, postId)
}

// IncrLike 增加点赞计数
func (i *InteractiveRepositoryImpl) IncrLike(ctx context.Context, postId uint, uid int64) error {
	return i.dao.InsertLikeInfo(ctx, dao.UserLikeBiz{
		BizID: postId,
		Uid:   uid,
	})
}

// DecrLike 减少点赞计数
func (i *InteractiveRepositoryImpl) DecrLike(ctx context.Context, postId uint, uid int64) error {
	return i.dao.DeleteLikeInfo(ctx, dao.UserLikeBiz{
		BizID: postId,
		Uid:   uid,
	})
}

// IncrCollectionItem 增加收藏计数
func (i *InteractiveRepositoryImpl) IncrCollectionItem(ctx context.Context, postId uint, uid int64) error {
	return i.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		BizID: postId,
		Uid:   uid,
	})
}

// DecrCollectionItem 减少收藏计数
func (i *InteractiveRepositoryImpl) DecrCollectionItem(ctx context.Context, postId uint, uid int64) error {
	return i.dao.DeleteCollectionBiz(ctx, dao.UserCollectionBiz{
		BizID: postId,
		Uid:   uid,
	})
}

// Get 获取互动信息
func (i *InteractiveRepositoryImpl) Get(ctx context.Context, postId uint) (domain.Interactive, error) {
	ic, err := i.dao.Get(ctx, postId)
	if err != nil {
		i.l.Error(PostGetInteractiveERROR, zap.Error(err))
		return domain.Interactive{}, err
	}

	return toDomain(ic), nil
}

// Liked 检查是否已点赞
func (i *InteractiveRepositoryImpl) Liked(ctx context.Context, postId uint, uid int64) (bool, error) {
	_, err := i.dao.GetLikeInfo(ctx, postId, uid)
	switch {
	case err == nil:
		// 如果没有错误，说明找到了点赞记录，返回true和重复点赞错误
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		// 如果错误是ErrRecordNotFound，说明没有找到点赞记录
		i.l.Error(PostGetLikedERROR, zap.Error(err))
		return false, nil
	default:
		// 如果是其他错误，返回false和错误信息
		return false, err
	}
}

// Collected 检查是否已收藏
func (i *InteractiveRepositoryImpl) Collected(ctx context.Context, postId uint, uid int64) (bool, error) {
	_, err := i.dao.GetCollectInfo(ctx, postId, uid)
	switch {
	case err == nil:
		// 如果没有错误，说明找到了收藏记录，返回true
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		// 如果错误是ErrRecordNotFound，说明没有找到收藏记录
		i.l.Error(PostGetCollectERROR, zap.Error(err))
		return false, nil
	default:
		// 如果是其他错误，返回false和错误信息
		return false, err
	}
}

// GetById 批量获取互动信息
func (i *InteractiveRepositoryImpl) GetById(ctx context.Context, postIds []uint) ([]domain.Interactive, error) {
	ics, err := i.dao.GetByIds(ctx, postIds)
	if err != nil {
		return make([]domain.Interactive, 0), err
	}
	result := make([]domain.Interactive, len(ics))
	for i, ic := range ics {
		result[i] = toDomain(ic)
	}

	return result, nil
}

func toDomain(ic dao.Interactive) domain.Interactive {
	return domain.Interactive{
		BizID:        ic.BizID,
		ReadCount:    ic.ReadCount,
		LikeCount:    ic.LikeCount,
		CollectCount: ic.CollectCount,
	}
}
