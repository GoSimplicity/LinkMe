package service

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"

	"go.uber.org/zap"
)

// InteractiveService 互动服务接口
type InteractiveService interface {
	Like(ctx context.Context, biz string, postId uint, uid int64) error                            // 点赞
	CancelLike(ctx context.Context, biz string, postId uint, uid int64) error                      // 取消点赞
	Collect(ctx context.Context, biz string, postId uint, cid, uid int64) error                    // 收藏
	CancelCollect(ctx context.Context, biz string, postId uint, cid, uid int64) error              // 取消收藏
	Get(ctx context.Context, biz string, postId uint) (domain.Interactive, error)                  // 获取互动信息
	GetByIds(ctx context.Context, biz string, postIds []uint) (map[uint]domain.Interactive, error) // 批量获取互动信息(热榜算法需要)
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

func (i *interactiveService) Like(ctx context.Context, biz string, postId uint, uid int64) error {
	liked, _ := i.repo.Liked(ctx, biz, postId, uid)
	// 如果已经点赞，则取消点赞
	if liked {
		return i.repo.DecrLike(ctx, biz, postId, uid)
	}

	return i.repo.IncrLike(ctx, biz, postId, uid)
}

func (i *interactiveService) CancelLike(ctx context.Context, biz string, postId uint, uid int64) error {
	return i.repo.DecrLike(ctx, biz, postId, uid)
}

func (i *interactiveService) Collect(ctx context.Context, biz string, postId uint, cid, uid int64) error {
	collected, _ := i.repo.Collected(ctx, biz, postId, uid)
	if collected {
		return i.repo.DecrCollectionItem(ctx, biz, postId, cid, uid)
	}

	return i.repo.IncrCollectionItem(ctx, biz, postId, cid, uid)
}

func (i *interactiveService) CancelCollect(ctx context.Context, biz string, postId uint, cid, uid int64) error {
	return i.repo.DecrCollectionItem(ctx, biz, postId, cid, uid)
}

func (i *interactiveService) Get(ctx context.Context, biz string, postId uint) (domain.Interactive, error) {
	di, err := i.repo.Get(ctx, biz, postId)
	if err != nil {
		return domain.Interactive{}, err
	}

	return di, err
}

func (i *interactiveService) GetByIds(ctx context.Context, biz string, postIds []uint) (map[uint]domain.Interactive, error) {
	dis, err := i.repo.GetById(ctx, biz, postIds)
	if err != nil {
		return nil, err
	}
	resultDis := make(map[uint]domain.Interactive)

	for _, interactive := range dis {
		resultDis[interactive.BizID] = interactive
	}

	return resultDis, err
}
