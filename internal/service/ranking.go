package service

import (
	"context"
	"errors"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/pkg/priorityqueue"
	"go.uber.org/zap"
	"math"
	"time"
)

type RankingService interface {
	TopN(ctx context.Context) error
	GetTopN(ctx context.Context) ([]domain.Post, error)
}

// rankingService 提供排名计算服务
type rankingService struct {
	interactiveService InteractiveService
	postRepository     repository.PostRepository
	rankingRepository  repository.RankingRepository
	l                  *zap.Logger
	batchSize          int                                                  // 每次分页处理的帖子数量
	rankSize           int                                                  // 要计算并返回的排名前 N 帖子的数量
	scoreFunc          func(likeCount int64, updatedTime time.Time) float64 // 用于计算帖子分数的函数，接受点赞数和更新时间作为参数，并返回计算后的分数
}

type Score struct {
	value float64     // 分数
	post  domain.Post // 帖子
}

func NewRankingService(interactiveService InteractiveService, postRepository repository.PostRepository, rankingRepository repository.RankingRepository, l *zap.Logger) RankingService {
	return &rankingService{
		interactiveService: interactiveService,
		postRepository:     postRepository,
		rankingRepository:  rankingRepository,
		l:                  l,
		batchSize:          100, // 默认设置分页处理的帖子数量为100
		rankSize:           100, // 默认设置要计算并返回的排名前 N 帖子的数量为100
		scoreFunc: func(likeCount int64, updatedTime time.Time) float64 { // 默认设置计算帖子分数的函数为点赞数减去1，除以更新时间加2的平方根
			duration := time.Since(updatedTime).Seconds()
			return float64(likeCount-1) / math.Pow(duration+2, 1.5)
		},
	}
}

// GetTopN 返回排名前 N 的帖子
func (b *rankingService) GetTopN(ctx context.Context) ([]domain.Post, error) {
	return b.rankingRepository.GetTopN(ctx)
}

// TopN 计算并替换排名前 N 的帖子
func (b *rankingService) TopN(ctx context.Context) error {
	posts, err := b.computeTopN(ctx)
	if err != nil {
		return err
	}
	return b.rankingRepository.ReplaceTopN(ctx, posts)
}

// computeTopN 计算排名前 N 的帖子
func (b *rankingService) computeTopN(ctx context.Context) ([]domain.Post, error) {
	offset := 0
	startTime := time.Now()
	// 将七天前的时间作为截止时间
	deadline := startTime.Add(-7 * 24 * time.Hour)
	// 初始化优先队列，比较函数根据 value 从大到小排序
	topNQueue := priorityqueue.NewPriorityQueue[Score](b.rankSize, func(a Score, b Score) bool {
		return a.value > b.value
	})
	for {
		// 分页处理
		posts, err := b.fetchPosts(ctx, offset)
		if err != nil {
			b.l.Error("fetch posts failed", zap.Error(err))
			return nil, err
		}
		// 如果没有更多帖子，跳出循环
		if len(posts) == 0 {
			break
		}
		// 获取每个帖子的交互数据
		interactions, err := b.fetchInteractions(ctx, posts)
		if err != nil {
			b.l.Error("fetch interactions failed", zap.Error(err))
			return nil, err
		}
		// 计算每个帖子的分数并加入优先队列
		for _, post := range posts {
			// 检查是否存在交互数据
			if interaction, ok := interactions[post.ID]; ok {
				b.l.Info("compute score", zap.Uint("postID", post.ID), zap.Int64("likeCount", interaction.LikeCount), zap.Time("updatedTime", post.UpdatedAt))
				// 使用 scoreFunc 计算分数
				score := b.scoreFunc(interaction.LikeCount, post.UpdatedAt)
				element := Score{
					value: score,
					post:  post,
				}
				// 将元素加入优先队列
				b.enqueueScore(topNQueue, element)
			}
		}
		// 更新偏移量
		offset += len(posts)
		// 检查是否需要终止循环
		if len(posts) < b.batchSize || posts[len(posts)-1].UpdatedAt.Before(deadline) {
			b.l.Info("compute topN done", zap.Int("offset", offset), zap.Int("batchSize", b.batchSize), zap.Int("rankSize", b.rankSize), zap.Duration("duration", time.Since(startTime)))
			break
		}
	}
	// 构建结果列表
	results := b.buildResults(topNQueue)
	return results, nil
}

// fetchPosts 获取完成分页后已发布的帖子
func (b *rankingService) fetchPosts(ctx context.Context, offset int) ([]domain.Post, error) {
	page := offset / b.batchSize
	size := int64(b.batchSize)
	pagination := domain.Pagination{
		Page: page,
		Size: &size,
	}
	return b.postRepository.ListPublishedPosts(ctx, pagination)
}

// fetchInteractions 获取每个帖子的交互数据
func (b *rankingService) fetchInteractions(ctx context.Context, posts []domain.Post) (map[uint]domain.Interactive, error) {
	// 创建帖子 ID 列表
	ids := make([]uint, len(posts))
	for i, post := range posts {
		ids[i] = post.ID
	}
	// 获取交互数据
	return b.interactiveService.GetByIds(ctx, "post", ids)
}

// enqueueScore 将元素加入优先队列，如果队列已满则进行替换
func (b *rankingService) enqueueScore(queue *priorityqueue.PriorityQueue[Score], element Score) {
	if err := queue.Enqueue(element); err != nil {
		if errors.Is(err, priorityqueue.ErrOutOfCapacity) {
			// 队列满了，取出最小元素
			minElement, _ := queue.Dequeue()
			if minElement.value < element.value {
				// 如果新元素分数高于最小元素，将新元素加入队列
				_ = queue.Enqueue(element)
			} else {
				// 否则，将最小元素重新加入队列
				_ = queue.Enqueue(minElement)
			}
		}
	}
}

// buildResults 从优先队列中构建结果列表
func (b *rankingService) buildResults(queue *priorityqueue.PriorityQueue[Score]) []domain.Post {
	results := make([]domain.Post, queue.Len())
	// 按分数从高到低依次取出元素
	for i := queue.Len() - 1; i >= 0; i-- {
		element, _ := queue.Dequeue()
		results[i] = element.post
	}
	return results
}
