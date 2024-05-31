package cache

import (
	"LinkMe/internal/domain"
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var (
	//go:embed lua/interactive_count.lua
	luaInteractiveCount string
	ErrKeyNotExist      = redis.Nil
)

const ReadCount = "read_count"
const LikeCount = "like_count"
const CollectCount = "collect_count"

type InteractiveCache interface {
	PostReadCountRecord(ctx context.Context, biz string, id int64) error         // 阅读计数
	PostLikeCountRecord(ctx context.Context, biz string, id int64) error         // 点赞计数
	DecrLikeCountRecord(ctx context.Context, biz string, id int64) error         // 取消点赞
	PostCollectCountRecord(ctx context.Context, biz string, id int64) error      // 收藏计数
	DecrCollectCountRecord(ctx context.Context, biz string, id int64) error      // 取消收藏
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)   // 获取互动信息
	Set(ctx context.Context, biz string, id int64, res domain.Interactive) error // 存储互动信息
}

type interactiveCache struct {
	client redis.Cmdable
}

func NewInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &interactiveCache{
		client: client,
	}
}

// Get 获取互动信息
func (i *interactiveCache) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	key := fmt.Sprintf("interactive:%s:%d", biz, id)
	res, err := i.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		return domain.Interactive{}, ErrKeyNotExist
	}
	var di domain.Interactive
	// 将获取到到计数转化为int64类型
	di.CollectCount, _ = strconv.ParseInt(res[CollectCount], 10, 64)
	di.LikeCount, _ = strconv.ParseInt(res[LikeCount], 10, 64)
	di.ReadCount, _ = strconv.ParseInt(res[ReadCount], 10, 64)
	return di, nil
}

// Set 存储互动信息
func (i *interactiveCache) Set(ctx context.Context,
	biz string, id int64,
	res domain.Interactive) error {
	// 设置键名
	key := fmt.Sprintf("interactive:%s:%d", biz, id)
	// 使用HashSet类型，写入收藏、阅读、点赞计数
	if err := i.client.HSet(ctx, key, CollectCount, res.CollectCount,
		ReadCount, res.ReadCount,
		LikeCount, res.LikeCount,
	).Err(); err != nil {
		return err
	}

	// 设置键的过期时间
	return i.client.Expire(ctx, key, time.Minute*15).Err()
}

// PostCollectCountRecord 收藏计数
func (i *interactiveCache) PostCollectCountRecord(ctx context.Context,
	biz string, id int64) error {
	key := fmt.Sprintf("interactive:%s:%d", biz, id)
	// 执行lua脚本，获取到计数加1
	return i.client.Eval(ctx, luaInteractiveCount, []string{key}, CollectCount, 1).Err()
}

// DecrCollectCountRecord 取消收藏
func (i *interactiveCache) DecrCollectCountRecord(ctx context.Context,
	biz string, id int64) error {
	key := fmt.Sprintf("interactive:%s:%d", biz, id)
	return i.client.Eval(ctx, luaInteractiveCount, []string{key}, CollectCount, -1).Err()
}

// PostLikeCountRecord 点赞计数
func (i *interactiveCache) PostLikeCountRecord(ctx context.Context,
	biz string, id int64) error {
	key := fmt.Sprintf("interactive:%s:%d", biz, id)
	return i.client.Eval(ctx, luaInteractiveCount, []string{key}, LikeCount, 1).Err()
}

// DecrLikeCountRecord 取消点赞
func (i *interactiveCache) DecrLikeCountRecord(ctx context.Context,
	biz string, id int64) error {
	key := fmt.Sprintf("interactive:%s:%d", biz, id)
	return i.client.Eval(ctx, luaInteractiveCount, []string{key}, LikeCount, -1).Err()
}

// PostReadCountRecord 阅读计数
func (i *interactiveCache) PostReadCountRecord(ctx context.Context,
	biz string, id int64) error {
	key := fmt.Sprintf("interactive:%s:%d", biz, id)
	return i.client.Eval(ctx, luaInteractiveCount, []string{key}, ReadCount, 1).Err()
}
