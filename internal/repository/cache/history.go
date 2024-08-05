package cache

import (
	"context"
	"encoding/json"
	"errors"
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
	// 从列表中获取元素
	values, err := h.client.LRange(ctx, key, 0, -1).Result()
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

// SetCache 将帖子设置在redis缓存中，并设置7天的过期时间
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
		return errors.New("无法获取锁")
	}
	// 当函数执行完毕，释放锁
	defer func() {
		if er := h.releaseLock(ctx, lockKey); er != nil {
			h.l.Error("释放锁失败", zap.Error(er))
		}
	}()
	// 获取缓存中的历史记录
	values, err := h.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		h.l.Error("获取缓存失败", zap.Error(err))
		return err
	}
	// 检查postId是否已经存在
	for _, v := range values {
		var existingHistory domain.History
		if er := json.Unmarshal([]byte(v), &existingHistory); er != nil {
			h.l.Error("反序列化历史记录失败", zap.Error(er))
			continue
		}
		if existingHistory.PostID == history.PostID {
			h.l.Info("历史记录已存在", zap.Uint("postId", history.PostID))
			return nil
		}
	}
	// 将帖子推送到 Redis 列表中
	if er := h.client.RPush(ctx, key, value).Err(); er != nil {
		h.l.Error("设置缓存失败", zap.Error(er))
		return er
	}
	// 设置过期时间为 7 天
	if er := h.client.Expire(ctx, key, 7*24*time.Hour).Err(); er != nil {
		h.l.Error("设置过期时间失败", zap.Error(er))
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
		return errors.New("无法获取锁")
	}
	// 当函数执行完毕，释放锁
	defer func() {
		if err := h.releaseLock(ctx, lockKey); err != nil {
			h.l.Error("释放锁失败", zap.Error(err))
		}
	}()
	// 获取缓存中的历史记录
	values, err := h.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		h.l.Error("获取缓存失败", zap.Error(err))
		return err
	}
	// 过滤数据，删除匹配的postID
	newValues := filterHistory(values, postId, h.l)
	// 清空原有的缓存记录
	if er := h.client.Del(ctx, key).Err(); er != nil {
		h.l.Error("删除缓存失败", zap.Error(er))
		return er
	}
	// 如果还有剩下的记录，重新推送到 Redis 列表中
	if len(newValues) > 0 {
		if er := h.client.RPush(ctx, key, newValues).Err(); er != nil {
			h.l.Error("重新设置缓存失败", zap.Error(er))
			return er
		}
	}
	h.l.Info("缓存删除成功", zap.String("key", key))
	return nil
}

// DeleteAllHistory 删除当前登录用户的全部历史记录缓存
func (h *historyCache) DeleteAllHistory(ctx context.Context, uid int64) error {
	pattern := fmt.Sprintf("linkme:post:history:%d", uid)
	var cursor uint64
	for {
		// 使用Scan命令分批获取匹配的键
		keys, nextCursor, err := h.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			h.l.Error("扫描历史记录键失败", zap.Error(err))
			return err
		}
		// 处理每组键
		for _, key := range keys {
			er := h.deleteKeyWithLock(ctx, key)
			if er != nil {
				// 记录日志但不中断所有删除操作
				h.l.Error("处理键失败", zap.String("key", key), zap.Error(er))
			}
		}
		// 更新游标
		cursor = nextCursor
		if cursor == 0 {
			break // 如果游标为0，说明扫描完成
		}
	}
	h.l.Info("所有历史记录已删除", zap.Int64("userID", uid))
	return nil
}

// deleteKeyWithLock 封装删除键操作，包括获取和释放锁
func (h *historyCache) deleteKeyWithLock(ctx context.Context, key string) error {
	lockKey := key + ":lock"
	lockTimeout := 5 * time.Second
	// 尝试获取锁
	locked, err := h.acquireLock(ctx, lockKey, lockTimeout)
	if err != nil {
		return fmt.Errorf("获取锁失败: %w", err)
	}
	if !locked {
		return fmt.Errorf("无法获取锁")
	}
	// 确保在此函数结束时释放锁
	defer func() {
		if er := h.releaseLock(ctx, lockKey); er != nil {
			h.l.Error("释放锁失败", zap.Error(er), zap.String("lockKey", lockKey))
		}
	}()
	// 执行删除操作
	if er := h.client.Del(ctx, key).Err(); er != nil {
		return fmt.Errorf("删除缓存失败: %w", er)
	}
	h.l.Info("缓存删除成功", zap.String("key", key))
	return nil
}

// filterHistory 过滤掉与指定postID匹配的历史记录
func filterHistory(values []string, postID uint, logger *zap.Logger) []string {
	var newValues []string
	for _, value := range values {
		var history domain.History
		if err := json.Unmarshal([]byte(value), &history); err != nil {
			logger.Error("反序列化历史记录失败", zap.Error(err))
			continue
		}
		if history.PostID != postID {
			newValues = append(newValues, value)
		}
	}
	return newValues
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
