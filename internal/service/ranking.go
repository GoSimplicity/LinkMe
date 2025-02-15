package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/job/interfaces"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/pkg/priorityqueue"
	"go.uber.org/zap"
)

var (
	ErrOffsetNegative = errors.New("偏移量不能为负数")
	ErrNilQueue       = errors.New("队列为空")
)

type RankingService interface {
	interfaces.RankingService
	GetTopN(ctx context.Context) ([]domain.Post, error)
	GetRankingConfig(ctx context.Context) (domain.RankingParameter, error)
	ResetRankingConfig(ctx context.Context, rankingParameter domain.RankingParameter) error
}

type Score struct {
	value float64
	post  domain.Post
}

type rankingService struct {
	interactiveRepository      repository.InteractiveRepository
	postRepository             repository.PostRepository
	rankingRepository          repository.RankingRepository
	rankingParameterRepository repository.RankingParameterRepository
	l                          *zap.Logger
	batchSize                  int
	rankSize                   int
	scoreFunc                  func(likeCount int64, ReadCount int64, CollectCount int64, updatedTime time.Time, rankingConfig domain.RankingParameter) float64
}

func NewRankingService(
	interactiveRepo repository.InteractiveRepository,
	postRepo repository.PostRepository,
	rankingRepo repository.RankingRepository,
	rankingParameterRepository repository.RankingParameterRepository,
	l *zap.Logger,
) RankingService {
	return &rankingService{
		interactiveRepository:      interactiveRepo,
		postRepository:             postRepo,
		rankingRepository:          rankingRepo,
		rankingParameterRepository: rankingParameterRepository,
		l:                          l,
		batchSize:                  100,
		rankSize:                   100,
		scoreFunc:                  calculateScore,
	}
}

// 核心业务方法
// TopN 计算并更新排名前 N 的帖子
func (rs *rankingService) TopN(ctx context.Context) error {
	posts, err := rs.computeTopN(ctx)
	if err != nil {
		return fmt.Errorf("计算前 N 名帖子失败: %w", err)
	}

	if len(posts) == 0 {
		rs.l.Info("没有需要排名的帖子")
		return nil
	}

	return rs.rankingRepository.ReplaceTopN(ctx, posts)
}

// GetTopN 获取排名前 N 的帖子
func (rs *rankingService) GetTopN(ctx context.Context) ([]domain.Post, error) {
	return rs.rankingRepository.GetTopN(ctx)
}

// 榜单配置相关方法
// GetRankingConfig 获取当前榜单配置
func (rs *rankingService) GetRankingConfig(ctx context.Context) (domain.RankingParameter, error) {
	return rs.rankingParameterRepository.FindLastParameter(ctx)
}

// ResetRankingConfig 重置榜单配置
func (rs *rankingService) ResetRankingConfig(ctx context.Context, rankingParameter domain.RankingParameter) error {
	_, err := rs.rankingParameterRepository.Insert(ctx, rankingParameter)
	if err != nil {
		return fmt.Errorf("重置榜单配置失败: %w", err)
	}
	return nil
}

// calculateScore 计算帖子得分
// 使用 HackerNews 算法的变体:score = likes / (t + 2)^1.5
// t 是发帖时间距现在的秒数,加2是为了避免除0
func calculateScore(likeCount int64, readCount int64, collectCount int64, updatedTime time.Time, rankingConfig domain.RankingParameter) float64 {
	// 处理默认值
	if rankingConfig.Alpha == 0 {
		rankingConfig.Alpha = 1.0
	}
	if rankingConfig.Beta == 0 {
		rankingConfig.Beta = 10.0
	}
	if rankingConfig.Gamma == 0 {
		rankingConfig.Gamma = 20.0
	}
	if rankingConfig.Lambda == 0 {
		rankingConfig.Lambda = 1.2
	}

	// 计算时间衰减因子
	duration := time.Since(updatedTime).Seconds()
	if duration < 0 {
		return 0
	}

	// 具体计算公式：
	// score = (Alpha * likeCount + Beta * collectCount + Gamma * readCount) / (duration + 2) ^ Lambda
	score := (rankingConfig.Alpha*float64(likeCount) +
		rankingConfig.Beta*float64(collectCount) +
		rankingConfig.Gamma*float64(readCount)) / math.Pow(duration+2, rankingConfig.Lambda)

	return score
}

// computeTopN 计算排名前 N 的帖子
func (rs *rankingService) computeTopN(ctx context.Context) ([]domain.Post, error) {
	topNQueue := priorityqueue.NewPriorityQueue[Score](rs.rankSize, func(a, b Score) bool {
		return a.value > b.value
	})

	offset := 0
	// 只处理最近7天的帖子
	deadline := time.Now().Add(-7 * 24 * time.Hour)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			posts, err := rs.processBatch(ctx, offset, topNQueue)
			if err != nil {
				return nil, err
			}

			// 处理完所有帖子或达到时间限制时退出
			if len(posts) == 0 || len(posts) < rs.batchSize || posts[len(posts)-1].UpdatedAt.Before(deadline) {
				return rs.buildResults(topNQueue), nil
			}

			offset += rs.batchSize
		}
	}
}

// processBatch 处理一批帖子
func (rs *rankingService) processBatch(ctx context.Context, offset int, queue *priorityqueue.PriorityQueue[Score]) ([]domain.Post, error) {
	if queue == nil {
		return nil, ErrNilQueue
	}

	posts, err := rs.fetchPosts(ctx, offset)
	if err != nil {
		return nil, fmt.Errorf("获取帖子失败: %w", err)
	}

	if len(posts) == 0 {
		return posts, nil
	}

	interactions, err := rs.fetchInteractions(ctx, posts)
	if err != nil {
		return nil, fmt.Errorf("获取交互数据失败: %w", err)
	}

	rankingConfigVal, err := rs.GetRankingConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取榜单配置失败: %w", err)
	}

	for _, post := range posts {
		if interaction, ok := interactions[post.ID]; ok {
			score := rs.scoreFunc(interaction.LikeCount, interaction.ReadCount, interaction.CollectCount, post.UpdatedAt, rankingConfigVal)
			rs.enqueueScore(queue, Score{value: score, post: post})
		}
	}

	return posts, nil
}

// fetchPosts 获取分页后的已发布帖子
func (rs *rankingService) fetchPosts(ctx context.Context, offset int) ([]domain.Post, error) {
	if offset < 0 {
		return nil, ErrOffsetNegative
	}
	size := int64(rs.batchSize)
	offsetInt64 := int64(offset)
	return rs.postRepository.ListPublishPosts(ctx, domain.Pagination{Size: &size, Offset: &offsetInt64}, 1) // 1 表示内部调用
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

	interactions, err := rs.interactiveRepository.GetById(ctx, ids)
	if err != nil {
		return nil, err
	}

	result := make(map[uint]domain.Interactive, len(interactions))
	for _, interaction := range interactions {
		result[interaction.BizID] = interaction
	}

	return result, nil
}

// enqueueScore 将分数加入优先队列
func (rs *rankingService) enqueueScore(queue *priorityqueue.PriorityQueue[Score], element Score) {
	if queue == nil {
		return
	}

	// 尝试直接入队
	if err := queue.Enqueue(element); err == nil {
		return
	}

	// 队列已满,需要与最小元素比较
	minElement, err := queue.Dequeue()
	if err != nil {
		rs.l.Error("出队列失败", zap.Error(err))
		return
	}

	// 选择较大的元素重新入队
	toEnqueue := element
	if minElement.value > element.value {
		toEnqueue = minElement
	}

	if err := queue.Enqueue(toEnqueue); err != nil {
		rs.l.Error("入队列失败", zap.Error(err))
	}
}

// buildResults 构建结果列表
func (rs *rankingService) buildResults(queue *priorityqueue.PriorityQueue[Score]) []domain.Post {
	if queue == nil || queue.Len() == 0 {
		return []domain.Post{}
	}

	results := make([]domain.Post, 0, queue.Len())
	for queue.Len() > 0 {
		element, err := queue.Dequeue()
		if err != nil {
			rs.l.Error("出队列失败", zap.Error(err))
			continue
		}
		results = append([]domain.Post{element.post}, results...)
	}

	return results
}
