package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type SMSCache interface {
	GetVCode(ctx context.Context, smsID, mobile string) string
	StoreVCode(ctx context.Context, smsID, mobile string, vCode string)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
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

func (s *smsCache) GetVCode(ctx context.Context, smsID, mobile string) string {
	key := getVCodeKey(smsID, mobile)
	vCode, _ := s.client.Get(ctx, key).Result()
	return vCode
}

func (s *smsCache) StoreVCode(ctx context.Context, smsID, mobile string, vCode string) {
	key := getVCodeKey(smsID, mobile)
	s.client.Set(ctx, key, vCode, time.Minute*10)
}

func (s *smsCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return s.client.SetNX(ctx, key, value, expiration)
}
