package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"LinkMe/pkg/priorityqueue"
	"LinkMe/pkg/slicetools"
	"context"
	"errors"
	"math"
	"time"
)

type RankingService interface {
	TopN(ctx context.Context) error
	GetTopN(ctx context.Context) ([]domain.Post, error)
}

type BatchRankingService struct {
	interactiveService InteractiveService
	postRepository     repository.PostRepository
	rankingRepository  repository.RankingRepository
	batchSize          int
	rankSize           int
	scoreFunc          func(likeCount int64, updatedTime time.Time) float64
}

func NewBatchRankingService(interactiveService InteractiveService, postRepository repository.PostRepository, rankingRepository repository.RankingRepository) RankingService {
	return &BatchRankingService{
		interactiveService: interactiveService,
		postRepository:     postRepository,
		rankingRepository:  rankingRepository,
		batchSize:          100,
		rankSize:           100,
		scoreFunc: func(likeCount int64, updatedTime time.Time) float64 {
			duration := time.Since(updatedTime).Seconds()
			return float64(likeCount-1) / math.Pow(duration+2, 1.5)
		},
	}
}

// GetTopN 返回排名前 N 的帖子
func (b *BatchRankingService) GetTopN(ctx context.Context) ([]domain.Post, error) {
	return b.rankingRepository.GetTopN(ctx)
}

// TopN 计算并替换排名前 N 的帖子
func (b *BatchRankingService) TopN(ctx context.Context) error {
	posts, err := b.computeTopN(ctx)
	if err != nil {
		return err
	}
	return b.rankingRepository.ReplaceTopN(ctx, posts)
}

// computeTopN 计算排名前 N 的帖子
func (b *BatchRankingService) computeTopN(ctx context.Context) ([]domain.Post, error) {
	offset := 0
	startTime := time.Now()
	deadline := startTime.Add(-7 * 24 * time.Hour)
	// Score 结构体包含分数和帖子
	type Score struct {
		value float64
		post  domain.Post
	}
	// 初始化优先队列
	topNQueue := priorityqueue.NewPriorityQueue[Score](b.rankSize,
		func(a Score, b Score) bool {
			return a.value > b.value
		})
	for {
		// 分页处理
		page := int64(offset / b.batchSize)
		size := int64(b.batchSize)
		pagination := domain.Pagination{
			Page: int(page),
			Size: &size,
			Uid:  0,
		}
		posts, err := b.postRepository.ListPublishedPosts(ctx, pagination)
		if err != nil {
			return nil, err
		}
		if len(posts) == 0 {
			break
		}
		// 获取帖子的 ID 列表
		ids := slicetools.Map(posts, func(_ int, post domain.Post) int64 {
			return post.ID
		})
		// 获取每个帖子的交互数据
		interactions, err := b.interactiveService.GetByIds(ctx, "Post", ids)
		if err != nil {
			return nil, err
		}
		// 计算每个帖子的分数并加入优先队列
		for _, post := range posts {
			interaction, found := interactions[post.ID]
			if !found {
				continue
			}
			score := b.scoreFunc(interaction.LikeCount, time.UnixMilli(post.UpdatedTime))
			element := Score{
				value: score,
				post:  post,
			}
			if er := topNQueue.Enqueue(element); er != nil {
				if errors.Is(er, priorityqueue.ErrOutOfCapacity) {
					minElement, _ := topNQueue.Dequeue()
					if minElement.value < score {
						_ = topNQueue.Enqueue(element)
					} else {
						_ = topNQueue.Enqueue(minElement)
					}
				} else {
					return nil, err
				}
			}
		}
		// 更新偏移量
		offset += len(posts)
		if len(posts) < int(size) || time.UnixMilli(posts[len(posts)-1].UpdatedTime).Before(deadline) {
			break
		}
	}
	// 构建结果列表
	results := make([]domain.Post, topNQueue.Len())
	for i := topNQueue.Len() - 1; i >= 0; i-- {
		element, _ := topNQueue.Dequeue()
		results[i] = element.post
	}
	return results, nil
}
