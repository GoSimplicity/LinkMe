package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/pkg/priorityqueue"
	"go.uber.org/zap"
)

type RankingService interface {
	TopN(ctx context.Context) error
	GetTopN(ctx context.Context) ([]domain.Post, error)
}

type Score struct {
	value float64     // 分数
	post  domain.Post // 帖子
}

type rankingService struct {
	interactiveService InteractiveService
	postRepository     repository.PostRepository
	rankingRepository  repository.RankingRepository
	l                  *zap.Logger
	batchSize          int                                                  // 每次分页处理的帖子数量
	rankSize           int                                                  // 要计算并返回的排名前 N 帖子的数量
	scoreFunc          func(likeCount int64, updatedTime time.Time) float64 // 用于计算帖子分数的函数
}

func NewRankingService(interactiveService InteractiveService, postRepository repository.PostRepository, rankingRepository repository.RankingRepository, l *zap.Logger) RankingService {
	return &rankingService{
		interactiveService: interactiveService,
		postRepository:     postRepository,
		rankingRepository:  rankingRepository,
		l:                  l,
		batchSize:          100,
		rankSize:           100,
		scoreFunc: func(likeCount int64, updatedTime time.Time) float64 {
			duration := time.Since(updatedTime).Seconds()
			if duration < 0 {
				return 0 // 防止未来时间
			}
			return float64(likeCount) / math.Pow(duration+2, 1.5)
		},
	}
}

// GetTopN 返回排名前 N 的帖子
func (rs *rankingService) GetTopN(ctx context.Context) ([]domain.Post, error) {
	return rs.rankingRepository.GetTopN(ctx)
}

// TopN 计算并替换排名前 N 的帖子
func (rs *rankingService) TopN(ctx context.Context) error {
	posts, err := rs.computeTopN(ctx)
	if err != nil {
		return fmt.Errorf("compute top N failed: %w", err)
	}
	return rs.rankingRepository.ReplaceTopN(ctx, posts)
}

// computeTopN 计算排名前 N 的帖子
func (rs *rankingService) computeTopN(ctx context.Context) ([]domain.Post, error) {
	offset := 0
	startTime := time.Now()
	deadline := startTime.Add(-7 * 24 * time.Hour)

	topNQueue := priorityqueue.NewPriorityQueue[Score](rs.rankSize, func(a, b Score) bool {
		return a.value > b.value
	})

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		posts, err := rs.fetchPosts(ctx, offset)
		if err != nil {
			rs.l.Error("获取帖子失败", zap.Error(err))
			return nil, fmt.Errorf("fetch posts failed: %w", err)
		}

		if len(posts) == 0 {
			break
		}

		interactions, err := rs.fetchInteractions(ctx, posts)
		if err != nil {
			rs.l.Error("获取交互数据失败", zap.Error(err))
			return nil, fmt.Errorf("fetch interactions failed: %w", err)
		}

		rs.processPostBatch(ctx, posts, interactions, topNQueue)

		offset += rs.batchSize

		if rs.shouldBreakLoop(posts, deadline) {
			rs.l.Info("计算完成",
				zap.Int("offset", offset),
				zap.Int("batchSize", rs.batchSize),
				zap.Int("rankSize", rs.rankSize),
				zap.Duration("duration", time.Since(startTime)),
			)
			break
		}
	}

	return rs.buildResults(topNQueue), nil
}

// processPostBatch 处理一批帖子
func (rs *rankingService) processPostBatch(ctx context.Context, posts []domain.Post, interactions map[uint]domain.Interactive, queue *priorityqueue.PriorityQueue[Score]) {
	for _, post := range posts {
		if interaction, ok := interactions[post.ID]; ok {
			score := rs.scoreFunc(interaction.LikeCount, post.UpdatedAt)
			rs.enqueueScore(queue, Score{value: score, post: post})
		}
	}
}

// shouldBreakLoop 判断是否应该结束循环
func (rs *rankingService) shouldBreakLoop(posts []domain.Post, deadline time.Time) bool {
	if len(posts) == 0 {
		return true
	}
	return len(posts) < rs.batchSize || posts[len(posts)-1].UpdatedAt.Before(deadline)
}

// fetchPosts 获取分页后的已发布帖子
func (rs *rankingService) fetchPosts(ctx context.Context, offset int) ([]domain.Post, error) {
	if offset < 0 {
		return nil, errors.New("offset cannot be negative")
	}
	page := offset / rs.batchSize
	size := int64(rs.batchSize)
	return rs.postRepository.ListPublishPosts(ctx, domain.Pagination{Page: page, Size: &size})
}

// fetchInteractions 获取帖子的交互数据
func (rs *rankingService) fetchInteractions(ctx context.Context, posts []domain.Post) (map[uint]domain.Interactive, error) {
	if len(posts) == 0 {
		return make(map[uint]domain.Interactive), nil
	}

	ids := make([]uint, len(posts))
	for i, post := range posts {
		ids[i] = post.ID
	}
	return rs.interactiveService.GetByIds(ctx, ids)
}

// enqueueScore 将分数加入优先队列
func (rs *rankingService) enqueueScore(queue *priorityqueue.PriorityQueue[Score], element Score) {
	if queue == nil {
		return
	}

	if err := queue.Enqueue(element); err != nil && errors.Is(err, priorityqueue.ErrOutOfCapacity) {
		minElement, err := queue.Dequeue()
		if err != nil {
			rs.l.Error("dequeue failed", zap.Error(err))
			return
		}
		if minElement.value < element.value {
			if err := queue.Enqueue(element); err != nil {
				rs.l.Error("enqueue failed", zap.Error(err))
			}
		} else {
			if err := queue.Enqueue(minElement); err != nil {
				rs.l.Error("enqueue failed", zap.Error(err))
			}
		}
	}
}

// buildResults 构建结果列表
func (rs *rankingService) buildResults(queue *priorityqueue.PriorityQueue[Score]) []domain.Post {
	if queue == nil {
		return nil
	}

	size := queue.Len()
	results := make([]domain.Post, size)
	for i := size - 1; i >= 0; i-- {
		element, err := queue.Dequeue()
		if err != nil {
			rs.l.Error("dequeue failed", zap.Error(err))
			continue
		}
		results[i] = element.post
	}
	return results
}
