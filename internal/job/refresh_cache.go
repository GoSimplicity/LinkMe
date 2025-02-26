package job

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

const (
	RefreshTypePub    = 1 // 刷新发布的帖子缓存
	RefreshTypeNormal = 2 // 刷新普通帖子缓存
	RefreshTypeAll    = 3 // 刷新所有类型缓存
)

type RefreshCacheTask struct {
	cache cache.PostCache
	l     *zap.Logger
}

func NewRefreshCacheTask(cache cache.PostCache, l *zap.Logger) *RefreshCacheTask {
	return &RefreshCacheTask{
		cache: cache,
		l:     l,
	}
}

type Payload struct {
	Key         string `json:"key"`
	PostId      uint   `json:"post_id"`
	RefreshType int    `json:"refresh_type"`
}

func (r *RefreshCacheTask) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p Payload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		r.l.Error("解析任务载荷失败", zap.Error(err))
		return fmt.Errorf("解析任务载荷失败: %w", err)
	}

	if p.RefreshType == 0 || p.Key == "" {
		return fmt.Errorf("无效的刷新参数: 类型=%d, 键=%s", p.RefreshType, p.Key)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	postIdStr := fmt.Sprint(p.PostId)
	logger := r.l.With(zap.Uint("post_id", p.PostId), zap.String("key", p.Key))

	// 根据刷新类型选择不同的删除策略
	switch p.RefreshType {
	case RefreshTypePub:
		return r.refreshPubCache(ctx, p.Key, postIdStr, logger)
	case RefreshTypeNormal:
		return r.refreshNormalCache(ctx, p.Key, postIdStr, logger)
	case RefreshTypeAll:
		return r.refreshAllCache(ctx, p.Key, postIdStr, logger)
	default:
		return fmt.Errorf("不支持的刷新类型: %d", p.RefreshType)
	}
}

// 刷新发布的帖子缓存
func (r *RefreshCacheTask) refreshPubCache(ctx context.Context, key, postId string, logger *zap.Logger) error {
	if err := r.cache.DelPubList(ctx, key); err != nil {
		logger.Error("删除发布帖子列表缓存失败", zap.Error(err))
		return fmt.Errorf("删除发布帖子列表缓存失败: %w", err)
	}

	if err := r.cache.DelPub(ctx, postId); err != nil {
		logger.Error("删除发布帖子缓存失败", zap.Error(err))
		return fmt.Errorf("删除发布帖子缓存失败: %w", err)
	}

	logger.Info("删除发布帖子缓存成功")
	return nil
}

// 刷新普通帖子缓存
func (r *RefreshCacheTask) refreshNormalCache(ctx context.Context, key, postId string, logger *zap.Logger) error {
	if err := r.cache.DelList(ctx, key); err != nil {
		logger.Error("删除普通帖子列表缓存失败", zap.Error(err))
		return fmt.Errorf("删除普通帖子列表缓存失败: %w", err)
	}

	if err := r.cache.Del(ctx, postId); err != nil {
		logger.Error("删除普通帖子缓存失败", zap.Error(err))
		return fmt.Errorf("删除普通帖子缓存失败: %w", err)
	}

	logger.Info("删除普通帖子缓存成功")
	return nil
}

// 刷新所有类型缓存
func (r *RefreshCacheTask) refreshAllCache(ctx context.Context, key, postId string, logger *zap.Logger) error {
	var errs []error

	if err := r.cache.DelList(ctx, key); err != nil {
		logger.Error("删除普通帖子列表缓存失败", zap.Error(err))
		errs = append(errs, fmt.Errorf("删除普通帖子列表缓存失败: %w", err))
	}

	if err := r.cache.Del(ctx, postId); err != nil {
		logger.Error("删除普通帖子缓存失败", zap.Error(err))
		errs = append(errs, fmt.Errorf("删除普通帖子缓存失败: %w", err))
	}

	if err := r.cache.DelPub(ctx, postId); err != nil {
		logger.Error("删除发布帖子缓存失败", zap.Error(err))
		errs = append(errs, fmt.Errorf("删除发布帖子缓存失败: %w", err))
	}

	if err := r.cache.DelPubList(ctx, key); err != nil {
		logger.Error("删除发布帖子列表缓存失败", zap.Error(err))
		errs = append(errs, fmt.Errorf("删除发布帖子列表缓存失败: %w", err))
	}

	if len(errs) > 0 {
		logger.Warn("部分缓存删除失败", zap.Int("error_count", len(errs)))
		return fmt.Errorf("部分缓存删除失败: %v", errs)
	}

	logger.Info("所有缓存删除成功")
	return nil
}
