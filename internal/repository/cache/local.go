package cache

import (
	"LinkMe/internal/domain"
	"LinkMe/pkg/atomicvalue"
	"context"
	"errors"
	"log"
	"time"
)

type RankingLocalCache interface {
	Set(ctx context.Context, arts []domain.Post) error
	Get(ctx context.Context) ([]domain.Post, error)
	ForceGet(ctx context.Context) ([]domain.Post, error)
}

type rankingLocalCache struct {
	topN       *atomicvalue.AtomicValue[[]domain.Post]
	ddl        *atomicvalue.AtomicValue[time.Time]
	expiration time.Duration
}

func NewRankingLocalCache(topN *atomicvalue.AtomicValue[[]domain.Post], ddl *atomicvalue.AtomicValue[time.Time], expiration time.Duration) RankingLocalCache {
	return &rankingLocalCache{
		topN:       topN,
		ddl:        ddl,
		expiration: expiration,
	}
}

func (r *rankingLocalCache) Set(ctx context.Context, arts []domain.Post) error {
	r.topN.Store(arts)
	r.ddl.Store(time.Now().Add(r.expiration))
	log.Printf("Cache set with %d posts, expires at %s", len(arts), r.ddl.Load())
	return nil
}

func (r *rankingLocalCache) Get(ctx context.Context) ([]domain.Post, error) {
	ddl := r.ddl.Load()
	arts := r.topN.Load()
	if len(arts) == 0 {
		log.Println("Local cache is empty.")
		return nil, errors.New("local cache is empty")
	}
	if ddl.Before(time.Now()) {
		log.Println("Local cache expired.")
		return nil, errors.New("local cache expired")
	}
	log.Printf("Cache hit with %d posts", len(arts))
	return arts, nil
}

func (r *rankingLocalCache) ForceGet(ctx context.Context) ([]domain.Post, error) {
	arts := r.topN.Load()
	if len(arts) == 0 {
		log.Println("Force get: local cache is empty.")
		return nil, errors.New("local cache is empty")
	}
	log.Printf("Force get: cache hit with %d posts", len(arts))
	return arts, nil
}
