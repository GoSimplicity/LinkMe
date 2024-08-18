package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"math/rand"
	"sync"
	"time"
)

type HistoryCache interface {
	SetCache(ctx context.Context, histories []domain.History) error
	GetCache(ctx context.Context, pagination domain.Pagination) ([]domain.History, error)
	DeleteOneCache(ctx context.Context, postId uint, uid int64) error
	DeleteAllHistory(ctx context.Context, uid int64) error
}

type historyCache struct {
	logger *zap.Logger
	client redis.Cmdable
	locks  sync.Map // 使用 sync.Map 存储本地锁
}

func NewHistoryCache(logger *zap.Logger, client redis.Cmdable) HistoryCache {
	return &historyCache{
		logger: logger,
		client: client,
	}
}

// GetCache 获取历史记录，并处理缓存穿透和缓存击穿
func (h *historyCache) GetCache(ctx context.Context, pagination domain.Pagination) ([]domain.History, error) {
	key := fmt.Sprintf("linkme:post:history:%d", pagination.Uid)
	threshold := time.Now().Add(-7 * 24 * time.Hour).Unix()

	// 本地锁，防止热点数据高并发场景下缓存击穿
	lockKey := key + ":lock"
	h.acquireLocalLock(lockKey)
	defer h.releaseLocalLock(lockKey)

	// 从缓存中获取数据
	values, err := h.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%d", threshold),
		Max:    "+inf",
		Offset: *pagination.Offset,
		Count:  *pagination.Size,
	}).Result()

	// 防止缓存穿透，设置空值缓存
	if errors.Is(err, redis.Nil) || len(values) == 0 {
		h.logger.Info("Cache miss, storing empty result to avoid cache penetration")

		err := h.client.Set(ctx, key+":empty", "1", 5*time.Minute).Err()
		if err != nil {
			h.logger.Error("Failed to set empty result cache", zap.Error(err))
			return nil, err
		} // 设置短暂过期时间的空值缓存，避免缓存穿透

		return nil, nil
	} else if err != nil {
		h.logger.Error("Failed to retrieve cache", zap.Error(err))
		return nil, err
	}

	// 防止缓存穿透的空值判断
	if h.client.Exists(ctx, key+":empty").Val() == 1 {
		h.logger.Info("Hit empty cache")
		return nil, nil
	}

	histories := make([]domain.History, 0, len(values))
	for _, v := range values {
		var history domain.History
		if err := json.Unmarshal([]byte(v), &history); err != nil {
			h.logger.Error("Failed to deserialize history record", zap.Error(err))
			continue
		}
		histories = append(histories, history)
	}

	return histories, nil
}

// SetCache 将数据写入缓存，并处理缓存雪崩
func (h *historyCache) SetCache(ctx context.Context, histories []domain.History) error {
	if len(histories) == 0 {
		return nil
	}

	key := fmt.Sprintf("linkme:post:history:%d", histories[0].AuthorID)
	lockKey := key + ":lock"

	// 本地锁，防止热点数据高并发场景下缓存击穿
	h.acquireLocalLock(lockKey)
	defer h.releaseLocalLock(lockKey)

	zAddArgs := make([]redis.Z, len(histories))
	for i, history := range histories {
		value, err := json.Marshal(history)
		if err != nil {
			h.logger.Error("Failed to serialize history record", zap.Error(err))
			return err
		}
		zAddArgs[i] = redis.Z{
			Score:  float64(time.Now().Unix()), // 使用时间戳作为分数存储在 ZSet 中
			Member: value,
		}
	}

	// 将数据写入缓存
	if err := h.client.ZAdd(ctx, key, zAddArgs...).Err(); err != nil {
		h.logger.Error("Failed to set cache", zap.Error(err))
		return err
	}

	// 设置随机过期时间，防止缓存雪崩
	expireDuration := 7*24*time.Hour + time.Duration(rand.Intn(3600))*time.Second // 通过随机增加过期时间防止缓存雪崩
	if err := h.client.Expire(ctx, key, expireDuration).Err(); err != nil {
		h.logger.Error("Failed to set expiration", zap.Error(err))
		return err
	}

	// 删除超过 7 天的记录
	threshold := time.Now().Add(-7 * 24 * time.Hour).Unix()
	if err := h.client.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", threshold)).Err(); err != nil {
		h.logger.Error("Failed to remove expired records", zap.Error(err))
		return err
	}

	h.logger.Info("Cache set successfully", zap.String("key", key))
	return nil
}

// DeleteOneCache 删除单条历史记录，防止缓存击穿
func (h *historyCache) DeleteOneCache(ctx context.Context, postId uint, uid int64) error {
	key := fmt.Sprintf("linkme:post:history:%d", uid)
	lockKey := key + ":lock"

	// 本地锁，防止热点数据高并发场景下缓存击穿
	h.acquireLocalLock(lockKey)
	defer h.releaseLocalLock(lockKey)

	values, err := h.client.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		h.logger.Error("Failed to retrieve cache", zap.Error(err))
		return err
	}

	for _, v := range values {
		var history domain.History
		if err := json.Unmarshal([]byte(v), &history); err != nil {
			h.logger.Error("Failed to deserialize history record", zap.Error(err))
			continue
		}
		if history.PostID == postId {
			if err := h.client.ZRem(ctx, key, v).Err(); err != nil {
				h.logger.Error("Failed to delete history record", zap.Error(err))
				return err
			}
			break
		}
	}

	h.logger.Info("History record deleted successfully", zap.String("key", key))
	return nil
}

// DeleteAllHistory 删除所有历史记录，并处理缓存击穿
func (h *historyCache) DeleteAllHistory(ctx context.Context, uid int64) error {
	key := fmt.Sprintf("linkme:post:history:%d", uid)
	lockKey := key + ":lock"

	// 本地锁，防止热点数据高并发场景下缓存击穿
	h.acquireLocalLock(lockKey)
	defer h.releaseLocalLock(lockKey)

	if err := h.client.Del(ctx, key).Err(); err != nil {
		h.logger.Error("Failed to delete all history records", zap.Error(err))
		return err
	}

	h.logger.Info("All history records deleted successfully", zap.Int64("userID", uid))
	return nil
}

// acquireLocalLock 获取本地锁，防止热点数据高并发场景下缓存击穿
func (h *historyCache) acquireLocalLock(lockKey string) {
	mu := &sync.Mutex{}
	actual, _ := h.locks.LoadOrStore(lockKey, mu) // 如果锁不存在，则存储新的锁
	mutex := actual.(*sync.Mutex)
	mutex.Lock()
}

// releaseLocalLock 释放本地锁
func (h *historyCache) releaseLocalLock(lockKey string) {
	if lock, ok := h.locks.Load(lockKey); ok {
		lock.(*sync.Mutex).Unlock()
	}
}

// acquireLock 尝试使用 Redis SetNX 获取分布式锁
func (h *historyCache) acquireLock(ctx context.Context, lockKey string, ttl time.Duration) (bool, error) {
	resp, err := h.client.SetNX(ctx, lockKey, "locked", ttl).Result() // 使用 SetNX 实现分布式锁
	if err != nil {
		return false, err
	}
	return resp, nil
}

// releaseLock 释放分布式锁，通过删除锁键来实现
func (h *historyCache) releaseLock(ctx context.Context, lockKey string) error {
	_, err := h.client.Del(ctx, lockKey).Result() // 通过删除锁键来释放分布式锁
	return err
}
