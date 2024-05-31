package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"

	"go.uber.org/zap"
)

// InteractiveService 互动服务接口
type InteractiveService interface {
	Like(ctx context.Context, biz string, id int64, uid int64) error                             // 点赞
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error                       // 取消点赞
	Collect(ctx context.Context, biz string, id, cid, uid int64) error                           //收藏
	CancelCollect(ctx context.Context, biz string, id, cid, uid int64) error                     // 取消收藏
	Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)        // 获取互动信息
	GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error) // 批量获取互动信息(热榜算法需要)
}

type interactiveService struct {
	repo repository.InteractiveRepository
	l    *zap.Logger
}

func NewInteractiveService(repo repository.InteractiveRepository, l *zap.Logger) InteractiveService {
	return &interactiveService{
		repo: repo,
		l:    l,
	}
}

func (i *interactiveService) Like(ctx context.Context, biz string, id int64, uid int64) error {
	liked, _ := i.repo.Liked(ctx, biz, id, uid)
	// 如果已经点赞，则取消点赞
	if liked {
		return i.repo.DecrLike(ctx, biz, id, uid)
	}
	return i.repo.IncrLike(ctx, biz, id, uid)
}

func (i *interactiveService) CancelLike(ctx context.Context, biz string, id int64, uid int64) error {
	return i.repo.DecrLike(ctx, biz, id, uid)
}

func (i *interactiveService) Collect(ctx context.Context, biz string, id, cid, uid int64) error {
	collected, _ := i.repo.Collected(ctx, biz, id, uid)
	if collected {
		return i.repo.DecrCollectionItem(ctx, biz, id, cid, uid)
	}
	return i.repo.IncrCollectionItem(ctx, biz, id, cid, uid)
}

func (i *interactiveService) CancelCollect(ctx context.Context, biz string, id, cid, uid int64) error {
	return i.repo.DecrCollectionItem(ctx, biz, id, cid, uid)
}

func (i *interactiveService) Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	di, err := i.repo.Get(ctx, biz, id)
	if err != nil {
		i.l.Error("get interactive filed", zap.Error(err))
		return domain.Interactive{}, err
	}
	return di, err
}

func (i *interactiveService) GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error) {
	dis, err := i.repo.GetById(ctx, biz, ids)
	if err != nil {
		i.l.Error("get interactions failed", zap.Error(err))
		return make(map[int64]domain.Interactive), err
	}
	result_dis := make(map[int64]domain.Interactive)

	for _, interactive := range dis {
		result_dis[interactive.BizID] = interactive
	}
	return result_dis, err
}
