package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type SMSCache interface {
	GetVCode(ctx context.Context, smsID, number string) (string, error)
	StoreVCode(ctx context.Context, smsID, number string, vCode string) error
	SetNX(ctx context.Context, number string, value interface{}, expiration time.Duration) (*redis.BoolCmd, error)
	Exist(ctx context.Context, number string) bool
	Count(ctx context.Context, number string) int
	IncrCnt(ctx context.Context, number string) error
	ReleaseLock(ctx context.Context, number string) error
}

type smsCache struct {
	client redis.Cmdable
}

func NewSMSCache(client redis.Cmdable) SMSCache {
	return &smsCache{
		client: client,
	}
}

const locked = "sms_locked"
const VCodeKey = "sms:%s:%s"
const LockedKey = "sms_locked:%s"
const CountKey = "sms:%s:%s"

func getVCodeKey(smsID, number string) string {
	return fmt.Sprintf(VCodeKey, smsID, number)
}
func getLockedKey(number string) string {
	return fmt.Sprintf(LockedKey, number)
}

func getCountKey(number string, now time.Time) string {
	return fmt.Sprintf(CountKey, number, now.Format("2006-01-02"))
}

// GetVCode 获取短信验证码
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

// StoreVCode 存储短信验证码
func (s *smsCache) StoreVCode(ctx context.Context, smsID string, number string, vCode string) error {
	key := getVCodeKey(smsID, number)
	return s.client.Set(ctx, key, vCode, time.Minute*10).Err()
}

// SetNX 设置一个不存在的key
func (s *smsCache) SetNX(ctx context.Context, number string, value interface{}, expiration time.Duration) (*redis.BoolCmd, error) {
	key := getLockedKey(number)
	result := s.client.SetNX(ctx, key, value, expiration)
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result, nil
}

// Exist 检查key是否存在
func (s *smsCache) Exist(ctx context.Context, number string) bool {
	key := getLockedKey(number)
	return s.client.Exists(ctx, key).Err() == nil
}

// Count 获取当天发送短信的次数
func (s *smsCache) Count(ctx context.Context, number string) int {
	key := getCountKey(number, time.Now())
	res, _ := strconv.ParseInt(s.client.Get(ctx, key).Val(), 10, 64)
	return int(res)
}

// IncrCnt 增加当天发送短信的次数
func (s *smsCache) IncrCnt(ctx context.Context, number string) error {
	now := time.Now()
	key := getCountKey(number, now)

	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	ttl := int(midnight.Sub(now).Seconds()) // 当天剩余的秒数

	// 使用 Lua 脚本来实现原子操作
	luaScript := `
		local current = redis.call('INCR', KEYS[1])
		if tonumber(current) == 1 then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end
    `

	// 执行 Lua 脚本
	_, err := s.client.Eval(ctx, luaScript, []string{key}, ttl).Result()
	return err
}

// ReleaseLock 释放 lock
func (s *smsCache) ReleaseLock(ctx context.Context, number string) error {
	luaScript := `
        local lock_key = KEYS[1]
        local lock_value = ARGV[1]

        if redis.call("GET", lock_key) == lock_value then
            redis.call("DEL", lock_key)
        end
        `

	if err := s.client.Eval(ctx, luaScript, []string{getLockedKey(number)}, locked).Err(); err != nil {
		return fmt.Errorf("failed to release lock")
	}
	return nil
}
