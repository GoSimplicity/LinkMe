package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type EmailCache interface {
	GetVCode(ctx context.Context, email string) (string, error)
	StoreVCode(ctx context.Context, email, vCode string) error
}

type emailCache struct {
	client redis.Cmdable
}

func NewEmailCache(client redis.Cmdable) EmailCache {
	return &emailCache{
		client: client,
	}
}

func (e emailCache) GetVCode(ctx context.Context, email string) (string, error) {
	return e.client.Get(ctx, genEmailKey(email)).Result()
}

func (e emailCache) StoreVCode(ctx context.Context, email, vCode string) error {
	return e.client.Set(ctx, genEmailKey(email), vCode, time.Duration(10)*time.Minute).Err()
}

func genEmailKey(email string) string {
	return fmt.Sprintf("linkme:email:%s", email)
}
