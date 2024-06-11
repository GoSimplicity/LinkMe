package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type SMSCache interface {
	GetVCode(ctx context.Context, smsID, number string) (string, error)
	StoreVCode(ctx context.Context, smsID, number string, vCode string) error
	SetNX(ctx context.Context, number string, value interface{}, expiration time.Duration) (*redis.BoolCmd, error)
	Exist(ctx context.Context, number string) bool
	Count(ctx context.Context, number string) int
	IncrCnt(ctx context.Context, number string) error
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
const LockedKey = "sms_locked:%s"

func getVCodeKey(smsID, number string) string {
	return fmt.Sprintf(VCodeKey, smsID, number)
}
func getLockedKey(number string) string {
	return fmt.Sprintf(LockedKey, number)
}

func (s *smsCache) GetVCode(ctx context.Context, smsID, number string) (string, error) {
	key := getVCodeKey(smsID, number)
	vCode, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil //键不存在或者过期了
		}
		return "", err
	}
	return vCode, nil
}

func (s *smsCache) StoreVCode(ctx context.Context, smsID string, number string, vCode string) error {
	key := getVCodeKey(smsID, number)
	return s.client.Set(ctx, key, vCode, time.Minute*10).Err()
}

func (s *smsCache) SetNX(ctx context.Context, number string, value interface{}, expiration time.Duration) (*redis.BoolCmd, error) {
	key := getLockedKey(number)
	result := s.client.SetNX(ctx, key, value, expiration)
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result, nil
}

func (s *smsCache) Exist(ctx context.Context, number string) bool {
	key := getLockedKey(number)
	return s.client.Exists(ctx, key).Err() == nil
}

func (s *smsCache) Count(ctx context.Context, number string) int {
	key := getLockedKey(number)
	res, _ := strconv.ParseInt(s.client.Get(ctx, key).Val(), 10, 64)
	return int(res)
}

func (s *smsCache) IncrCnt(ctx context.Context, number string) error {
	key := getLockedKey(number)

	// 获取当前时间和当天结束的时间
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	ttl := int(midnight.Sub(now).Seconds()) // 当天剩余的秒数

	// 使用 Lua 脚本来实现原子操作
	luaScript := `
    local current = redis.call('INCR', KEYS[1])
    if current == 1 then
        redis.call('EXPIRE', KEYS[1], ARGV[1])
    end
    return current
    `

	// 执行 Lua 脚本
	_, err := s.client.Eval(ctx, luaScript, []string{key}, ttl).Result()
	if err != nil {
		return err
	}

	return nil
}
