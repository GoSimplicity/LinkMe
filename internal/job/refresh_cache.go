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
	Key    string `json:"key"`
	PostId uint   `json:"post_id"`
}

func (r *RefreshCacheTask) ProcessTask(ctx context.Context, t *asynq.Task) error {
	// 解析任务载荷
	var p Payload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		r.l.Error("解析任务载荷失败", zap.Error(err))
		return fmt.Errorf("解析任务载荷失败: %w", err)
	}

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 删除帖子列表缓存
	if err := r.cache.DelPubList(ctx, p.Key); err != nil {
		r.l.Error("删除帖子列表缓存失败",
			zap.Error(err),
			zap.String("key", p.Key))
		return fmt.Errorf("删除帖子列表缓存失败: %w", err)
	}

	// 如果存在帖子ID,则删除对应的已发布帖子缓存
	if p.PostId > 0 {
		if err := r.cache.DelPub(ctx, p.Key); err != nil {
			r.l.Error("删除已发布帖子缓存失败",
				zap.Error(err),
				zap.String("key", p.Key),
				zap.Uint("post_id", p.PostId))
			return fmt.Errorf("删除已发布帖子缓存失败: %w", err)
		}
	}

	r.l.Info("成功删除帖子缓存",
		zap.String("key", p.Key),
		zap.Uint("post_id", p.PostId))

	return nil
}
