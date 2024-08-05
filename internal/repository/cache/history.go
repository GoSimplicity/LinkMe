package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type HistoryCache interface {
	SetCache(ctx context.Context, history domain.History) error
	GetCache(ctx context.Context, pagination domain.Pagination) ([]domain.History, error)
	DeleteOneCache(ctx context.Context, postId uint, uid int64) error
	DeleteAllHistory(ctx context.Context, uid int64) error
}

type historyCache struct {
	l      *zap.Logger
	client redis.Cmdable
}

func NewHistoryCache(l *zap.Logger, client redis.Cmdable) HistoryCache {
	return &historyCache{
		l:      l,
		client: client,
	}
}

// GetCache 根据分页信息从Redis缓存中获取历史记录列表
func (h *historyCache) GetCache(ctx context.Context, pagination domain.Pagination) ([]domain.History, error) {
	key := fmt.Sprintf("linkme:post:history:%d", pagination.Uid)
	// 获取当前时间戳，用于过滤超过7天的记录
	threshold := time.Now().Add(-7 * 24 * time.Hour).Unix()
	// 从zset中获取最近7天的历史记录
	values, err := h.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%d", threshold),
		Max:    "+inf",
		Offset: *pagination.Offset,
		Count:  *pagination.Size,
	}).Result()
	if err != nil {
		h.l.Error("获取缓存失败", zap.Error(err))
		return nil, err
	}
	var histories []domain.History
	for _, v := range values {
		var history domain.History
		if er := json.Unmarshal([]byte(v), &history); er != nil {
			h.l.Error("反序列化历史记录失败", zap.Error(er))
			continue
		}
		histories = append(histories, history)
	}
	return histories, nil
}

// SetCache 将帖子设置在redis缓存中 并设置7天的过期时间
func (h *historyCache) SetCache(ctx context.Context, history domain.History) error {
	key := fmt.Sprintf("linkme:post:history:%d", history.AuthorID)
	value, err := json.Marshal(history)
	if err != nil {
		h.l.Error("序列化帖子失败", zap.Error(err))
		return err
	}
	// 定义锁键和锁超时时间
	lockKey := key + ":lock"
	lockTimeout := 5 * time.Second
	// 尝试获取锁
	locked, err := h.acquireLock(ctx, lockKey, lockTimeout)
	if err != nil {
		h.l.Error("获取锁失败", zap.Error(err))
		return err
	}
	if !locked {
		return fmt.Errorf("无法获取锁")
	}
	// 当函数执行完毕，释放锁
	defer func() {
		if er := h.releaseLock(ctx, lockKey); er != nil {
			h.l.Error("释放锁失败", zap.Error(er))
		}
	}()
	// 将帖子记录添加到zset中，使用时间戳作为分数
	if er := h.client.ZAdd(ctx, key, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: value,
	}).Err(); er != nil {
		h.l.Error("设置缓存失败", zap.Error(er))
		return er
	}
	// 删除超过7天的记录
	threshold := time.Now().Add(-7 * 24 * time.Hour).Unix()
	if er := h.client.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", threshold)).Err(); er != nil {
		h.l.Error("删除过期记录失败", zap.Error(er))
		return er
	}
	h.l.Info("缓存设置成功", zap.String("key", key))
	return nil
}

// DeleteOneCache 删除缓存中的一条历史记录，前提是post被标记为删除
func (h *historyCache) DeleteOneCache(ctx context.Context, postId uint, uid int64) error {
	key := fmt.Sprintf("linkme:post:history:%d", uid)
	lockKey := key + ":lock"
	lockTimeout := 5 * time.Second
	// 尝试获取锁
	if locked, err := h.acquireLock(ctx, lockKey, lockTimeout); err != nil {
		h.l.Error("获取锁失败", zap.Error(err))
		return err
	} else if !locked {
		return fmt.Errorf("无法获取锁")
	}
	// 当函数执行完毕，释放锁
	defer func() {
		if err := h.releaseLock(ctx, lockKey); err != nil {
			h.l.Error("释放锁失败", zap.Error(err))
		}
	}()
	// 获取缓存中的历史记录
	values, err := h.client.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		h.l.Error("获取缓存失败", zap.Error(err))
		return err
	}
	// 查找并删除匹配的postID
	for _, v := range values {
		var history domain.History
		if er := json.Unmarshal([]byte(v), &history); er != nil {
			h.l.Error("反序列化历史记录失败", zap.Error(er))
			continue
		}
		if history.PostID == postId {
			if er := h.client.ZRem(ctx, key, v).Err(); er != nil {
				h.l.Error("删除历史记录失败", zap.Error(er))
				return er
			}
			break
		}
	}
	h.l.Info("缓存删除成功", zap.String("key", key))
	return nil
}

// DeleteAllHistory 删除当前登录用户的全部历史记录缓存
func (h *historyCache) DeleteAllHistory(ctx context.Context, uid int64) error {
	key := fmt.Sprintf("linkme:post:history:%d", uid)
	lockKey := key + ":lock"
	lockTimeout := 5 * time.Second
	// 尝试获取锁
	if locked, err := h.acquireLock(ctx, lockKey, lockTimeout); err != nil {
		h.l.Error("获取锁失败", zap.Error(err))
		return err
	} else if !locked {
		return fmt.Errorf("无法获取锁")
	}
	// 当函数执行完毕，释放锁
	defer func() {
		if err := h.releaseLock(ctx, lockKey); err != nil {
			h.l.Error("释放锁失败", zap.Error(err))
		}
	}()
	// 删除zset
	if err := h.client.Del(ctx, key).Err(); err != nil {
		h.l.Error("删除历史记录失败", zap.Error(err))
		return err
	}
	h.l.Info("所有历史记录已删除", zap.Int64("userID", uid))
	return nil
}

// acquireLock 尝试使用 Redis SetNX 获取分布式锁
func (h *historyCache) acquireLock(ctx context.Context, lockKey string, ttl time.Duration) (bool, error) {
	resp, err := h.client.SetNX(ctx, lockKey, "locked", ttl).Result()
	if err != nil {
		return false, err
	}
	return resp, nil
}

// releaseLock 释放分布式锁，通过删除锁键来实现
func (h *historyCache) releaseLock(ctx context.Context, lockKey string) error {
	_, err := h.client.Del(ctx, lockKey).Result()
	return err
}
