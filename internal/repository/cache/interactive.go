package cache

import (
	"context"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"time"
)

var (
	ErrKeyNotExist = redis.Nil
)

const (
	ReadCount    = "read_count"
	LikeCount    = "like_count"
	CollectCount = "collect_count"
)

type InteractiveCache interface {
	PostReadCountRecord(ctx context.Context, biz string, postId uint) error         // 阅读计数
	PostLikeCountRecord(ctx context.Context, biz string, postId uint) error         // 点赞计数
	DecrLikeCountRecord(ctx context.Context, biz string, postId uint) error         // 取消点赞
	PostCollectCountRecord(ctx context.Context, biz string, postId uint) error      // 收藏计数
	DecrCollectCountRecord(ctx context.Context, biz string, postId uint) error      // 取消收藏
	Get(ctx context.Context, biz string, postId uint) (domain.Interactive, error)   // 获取互动信息
	Set(ctx context.Context, biz string, postId uint, res domain.Interactive) error // 存储互动信息
}

type interactiveCache struct {
	client       redis.Cmdable
	locker       *redislock.Client
	incrByScript *redis.Script
}

func NewInteractiveCache(client redis.Cmdable) InteractiveCache {
	locker := redislock.New(client)

	// 将 Lua 脚本加载到 Redis 中，并保存其 SHA1 值以便后续调用
	incrByScript := redis.NewScript(`
local key = KEYS[1]
local field = ARGV[1]
local increment = tonumber(ARGV[2])
local value = redis.call("HINCRBY", key, field, increment)
return value
`)

	return &interactiveCache{
		client:       client,
		locker:       locker,
		incrByScript: incrByScript,
	}
}

// Get 获取互动信息
func (i *interactiveCache) Get(ctx context.Context, biz string, postId uint) (domain.Interactive, error) {
	key := fmt.Sprintf("interactive:%s:%d", biz, postId)
	res, err := i.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		return domain.Interactive{}, ErrKeyNotExist
	}
	var di domain.Interactive
	di.CollectCount, _ = strconv.ParseInt(res[CollectCount], 10, 64)
	di.LikeCount, _ = strconv.ParseInt(res[LikeCount], 10, 64)
	di.ReadCount, _ = strconv.ParseInt(res[ReadCount], 10, 64)
	return di, nil
}

// Set 存储互动信息
func (i *interactiveCache) Set(ctx context.Context, biz string, postId uint, res domain.Interactive) error {
	key := fmt.Sprintf("interactive:%s:%d", biz, postId)
	err := i.client.HSet(ctx, key, CollectCount, res.CollectCount,
		ReadCount, res.ReadCount,
		LikeCount, res.LikeCount,
	).Err()
	if err != nil {
		return err
	}

	// 设置随机化的键的过期时间，防止缓存雪崩
	expiration := time.Hour*24 + time.Duration(rand.Intn(3600))*time.Second
	return i.client.Expire(ctx, key, expiration).Err()
}

// PostCollectCountRecord 收藏计数
func (i *interactiveCache) PostCollectCountRecord(ctx context.Context, biz string, postId uint) error {
	key := fmt.Sprintf("interactive:%s:%d", biz, postId)
	_, err := i.incrByScript.Run(ctx, i.client, []string{key}, CollectCount, 1).Result()
	return err
}

// DecrCollectCountRecord 取消收藏
func (i *interactiveCache) DecrCollectCountRecord(ctx context.Context, biz string, postId uint) error {
	key := fmt.Sprintf("interactive:%s:%d", biz, postId)
	_, err := i.incrByScript.Run(ctx, i.client, []string{key}, CollectCount, -1).Result()
	return err
}

// PostLikeCountRecord 点赞计数
func (i *interactiveCache) PostLikeCountRecord(ctx context.Context, biz string, postId uint) error {
	key := fmt.Sprintf("interactive:%s:%d", biz, postId)
	_, err := i.incrByScript.Run(ctx, i.client, []string{key}, LikeCount, 1).Result()
	return err
}

// DecrLikeCountRecord 取消点赞
func (i *interactiveCache) DecrLikeCountRecord(ctx context.Context, biz string, postId uint) error {
	key := fmt.Sprintf("interactive:%s:%d", biz, postId)
	_, err := i.incrByScript.Run(ctx, i.client, []string{key}, LikeCount, -1).Result()
	return err
}

// PostReadCountRecord 阅读计数
func (i *interactiveCache) PostReadCountRecord(ctx context.Context, biz string, postId uint) error {
	key := fmt.Sprintf("interactive:%s:%d", biz, postId)
	_, err := i.incrByScript.Run(ctx, i.client, []string{key}, ReadCount, 1).Result()
	return err
}
