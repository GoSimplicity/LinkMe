package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RankingRedisCache interface {
	Set(ctx context.Context, arts []domain.Post) error
	Get(ctx context.Context) ([]domain.Post, error)
}

type rankingCache struct {
	client     redis.Cmdable
	key        string
	expiration time.Duration
}

func NewRankingRedisCache(client redis.Cmdable) RankingRedisCache {
	return &rankingCache{
		client:     client,
		key:        "ranking:top_n",
		expiration: 3 * time.Minute,
	}
}

func (r *rankingCache) Set(ctx context.Context, arts []domain.Post) error {
	for i := range arts {
		arts[i].Content = arts[i].Abstract()
	}
	val, err := json.Marshal(arts)
	if err != nil {
		log.Printf("Error marshalling posts: %v", err)
		return err
	}
	if er := r.client.Set(ctx, r.key, val, r.expiration).Err(); er != nil {
		log.Printf("Error setting cache: %v", er)
		return er
	}
	return nil
}

func (r *rankingCache) Get(ctx context.Context) ([]domain.Post, error) {
	val, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Printf("Cache miss for key: %s", r.key)
			return nil, nil // Cache miss is not an error
		}
		log.Printf("Error getting cache: %v", err)
		return nil, err
	}
	var res []domain.Post
	if er := json.Unmarshal(val, &res); er != nil {
		log.Printf("Error unmarshalling posts: %v", er)
		return nil, er
	}
	return res, nil
}
