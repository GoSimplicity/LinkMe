package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// 历史记录保留时间为7天
	historyRetentionDays = 7 * 24 * time.Hour
	// 空值缓存过期时间为5分钟
	emptyCacheTTL = 5 * time.Minute
	// 缓存键前缀
	historyKeyPrefix = "linkme:post:history:"
)

type HistoryCache interface {
	SetCache(ctx context.Context, history domain.History) error
	GetCache(ctx context.Context, pagination domain.Pagination) ([]domain.History, error)
	DeleteOneCache(ctx context.Context, postId uint, uid int64) error
	DeleteAllHistory(ctx context.Context, uid int64) error
}

type historyCache struct {
	logger *zap.Logger
	client redis.Cmdable
	locks  sync.Map
}

func NewHistoryCache(logger *zap.Logger, client redis.Cmdable) HistoryCache {
	return &historyCache{
		logger: logger,
		client: client,
	}
}

// GetCache 获取历史记录
func (h *historyCache) GetCache(ctx context.Context, pagination domain.Pagination) ([]domain.History, error) {
	key := fmt.Sprintf("%s%d", historyKeyPrefix, pagination.Uid)
	threshold := time.Now().Add(-historyRetentionDays).Unix()

	var histories []domain.History
	var err error

	// 使用本地锁防止缓存击穿
	err = h.withLocalLock(key+":lock", func() error {
		// 从缓存中获取数据
		values, err := h.client.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
			Min:    fmt.Sprintf("%d", threshold),
			Max:    "+inf",
			Offset: *pagination.Offset,
			Count:  *pagination.Size,
		}).Result()

		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil
			}
			h.logger.Error("获取历史记录失败", zap.Error(err))
			return err
		}

		if len(values) == 0 {
			return nil
		}

		// 反序列化历史记录
		histories = make([]domain.History, 0, len(values))
		for _, v := range values {
			var history domain.History
			if err := json.Unmarshal([]byte(v), &history); err != nil {
				h.logger.Error("反序列化历史记录失败", zap.Error(err))
				continue
			}
			histories = append(histories, history)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return histories, nil
}

func (h *historyCache) SetCache(ctx context.Context, history domain.History) error {
	key := fmt.Sprintf("%s%d", historyKeyPrefix, history.Uid)

	return h.withLocalLock(key+":lock", func() error {
		value, err := json.Marshal(history)
		if err != nil {
			h.logger.Error("序列化历史记录失败", zap.Error(err))
			return err
		}

		// 使用管道执行多个操作
		pipe := h.client.Pipeline()
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: value,
		})
		pipe.Expire(ctx, key, historyRetentionDays+time.Duration(rand.Intn(3600))*time.Second)
		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", time.Now().Add(-historyRetentionDays).Unix()))

		if _, err := pipe.Exec(ctx); err != nil {
			h.logger.Error("缓存操作失败", zap.Error(err))
			return err
		}

		h.logger.Info("历史记录缓存成功", zap.String("key", key))
		return nil
	})
}

func (h *historyCache) DeleteOneCache(ctx context.Context, postId uint, uid int64) error {
	key := fmt.Sprintf("%s%d", historyKeyPrefix, uid)

	return h.withLocalLock(key+":lock", func() error {
		values, err := h.client.ZRange(ctx, key, 0, -1).Result()
		if err != nil {
			h.logger.Error("获取历史记录失败", zap.Error(err))
			return err
		}

		for _, v := range values {
			var history domain.History
			if err := json.Unmarshal([]byte(v), &history); err != nil {
				h.logger.Error("反序列化历史记录失败", zap.Error(err))
				continue
			}
			if history.PostID == postId {
				if err := h.client.ZRem(ctx, key, v).Err(); err != nil {
					h.logger.Error("删除历史记录失败", zap.Error(err))
					return err
				}
				break
			}
		}

		h.logger.Info("历史记录删除成功", zap.String("key", key))
		return nil
	})
}

func (h *historyCache) DeleteAllHistory(ctx context.Context, uid int64) error {
	key := fmt.Sprintf("%s%d", historyKeyPrefix, uid)

	return h.withLocalLock(key+":lock", func() error {
		if err := h.client.Del(ctx, key).Err(); err != nil {
			h.logger.Error("删除所有历史记录失败", zap.Error(err))
			return err
		}

		h.logger.Info("所有历史记录删除成功", zap.Int64("uid", uid))
		return nil
	})
}

// 工具函数

func (h *historyCache) withLocalLock(lockKey string, fn func() error) error {
	mu := &sync.Mutex{}
	actual, _ := h.locks.LoadOrStore(lockKey, mu)
	mutex := actual.(*sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()
	return fn()
}

func (h *historyCache) handleEmptyOrError(ctx context.Context, key string, values []string, err error) bool {
	if errors.Is(err, redis.Nil) || len(values) == 0 {
		h.logger.Info("缓存未命中,设置空值缓存")
		if err := h.client.Set(ctx, key+":empty", "1", emptyCacheTTL).Err(); err != nil {
			h.logger.Error("设置空值缓存失败", zap.Error(err))
		}
		return true
	}

	if err != nil {
		h.logger.Error("获取缓存失败", zap.Error(err))
		return true
	}

	if h.client.Exists(ctx, key+":empty").Val() == 1 {
		h.logger.Info("命中空值缓存")
		return true
	}

	return false
}
