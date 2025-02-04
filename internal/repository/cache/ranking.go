package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	ErrMarshalPost   = errors.New("序列化帖子失败")
	ErrUnmarshalPost = errors.New("反序列化帖子失败")
	ErrSetCache      = errors.New("设置缓存失败")
	ErrGetCache      = errors.New("获取缓存失败")
)

type RankingRedisCache interface {
	Set(ctx context.Context, posts []domain.Post) error
	Get(ctx context.Context) ([]domain.Post, error)
}

type rankingCache struct {
	client     redis.Cmdable
	key        string
	expiration time.Duration
	logger     *zap.Logger
}

func NewRankingRedisCache(client redis.Cmdable, logger *zap.Logger) RankingRedisCache {
	return &rankingCache{
		client:     client,
		key:        "ranking:top_n",
		expiration: 3 * time.Minute,
		logger:     logger,
	}
}

func (r *rankingCache) Set(ctx context.Context, posts []domain.Post) error {
	// 预处理帖子内容
	for i := range posts {
		posts[i].Content = posts[i].Abstract()
	}

	// 序列化帖子数据
	val, err := json.Marshal(posts)
	if err != nil {
		r.logger.Error("序列化帖子失败", zap.Error(err))
		return ErrMarshalPost
	}

	// 设置缓存
	if err := r.client.Set(ctx, r.key, val, r.expiration).Err(); err != nil {
		r.logger.Error("设置缓存失败",
			zap.String("key", r.key),
			zap.Error(err))
		return ErrSetCache
	}

	r.logger.Info("缓存设置成功",
		zap.String("key", r.key),
		zap.Int("post_count", len(posts)))
	return nil
}

func (r *rankingCache) Get(ctx context.Context) ([]domain.Post, error) {
	// 获取缓存数据
	val, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.logger.Info("缓存未命中", zap.String("key", r.key))
			return nil, nil
		}
		r.logger.Error("获取缓存失败",
			zap.String("key", r.key),
			zap.Error(err))
		return nil, ErrGetCache
	}

	// 反序列化帖子数据
	var posts []domain.Post
	if err := json.Unmarshal(val, &posts); err != nil {
		r.logger.Error("反序列化帖子失败", zap.Error(err))
		return nil, ErrUnmarshalPost
	}

	r.logger.Info("缓存获取成功",
		zap.String("key", r.key),
		zap.Int("post_count", len(posts)))

	return posts, nil
}
