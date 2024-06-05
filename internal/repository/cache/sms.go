package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type SMSCache interface {
	GetVCode(ctx context.Context, smsID, mobile string) (string, error)
	StoreVCode(ctx context.Context, smsID, mobile string, vCode string) error
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (*redis.BoolCmd, error)
}

type smsCache struct {
	client redis.Cmdable
}

func NewSMSCache(client redis.Cmdable) SMSCache {
	return &smsCache{
		client: client,
	}
}

const VCodeKey = "sms:%s:%s"

func getVCodeKey(smsID, mobile string) string {
	return fmt.Sprintf(VCodeKey, smsID, mobile)
}

func (s *smsCache) GetVCode(ctx context.Context, smsID, mobile string) (string, error) {
	key := getVCodeKey(smsID, mobile)
	vCode, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil //键不存在或者过期了
		}
		return "", err
	}
	return vCode, nil
}

func (s *smsCache) StoreVCode(ctx context.Context, smsID, mobile string, vCode string) error {
	key := getVCodeKey(smsID, mobile)
	return s.client.Set(ctx, key, vCode, time.Minute*10).Err()
}

func (s *smsCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (*redis.BoolCmd, error) {
	result := s.client.SetNX(ctx, key, value, expiration)
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result, nil
}
