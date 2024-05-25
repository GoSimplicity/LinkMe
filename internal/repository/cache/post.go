package cache

import (
	"LinkMe/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type PostCache interface {
	GetFirstPage(ctx context.Context, id int64) ([]domain.Post, error)       // 获取个人用户的第一页帖子缓存
	GetPubFirstPage(ctx context.Context, id int64) ([]domain.Post, error)    // 获取公开用户的第一页帖子缓存
	SetFirstPage(ctx context.Context, id int64, post []domain.Post) error    // 设置个人用户的第一页帖子缓存
	SetPubFirstPage(ctx context.Context, id int64, post []domain.Post) error // 设置公开用户的第一页帖子缓存
	DelFirstPage(ctx context.Context, id int64) error                        // 删除用户的第一页帖子缓存
	GetDetail(ctx context.Context, id int64) (domain.Post, error)            // 根据ID获取一个帖子详情缓存
	SetDetail(ctx context.Context, post domain.Post) error                   // 设置一个帖子详情缓存
	GetPubDetail(ctx context.Context, id int64) (domain.Post, error)         // 根据ID获取一个已发布的帖子详情缓存
	SetPubDetail(ctx context.Context, post domain.Post) error                // 设置一个已发布的帖子详情缓存
}

type postCache struct {
	cmd redis.Cmdable
	l   *zap.Logger
}

func NewPostCache(cmd redis.Cmdable, l *zap.Logger) PostCache {
	return &postCache{
		cmd: cmd,
		l:   l,
	}
}

// GetFirstPage 获取个人第一页帖子摘要
func (p *postCache) GetFirstPage(ctx context.Context, id int64) ([]domain.Post, error) {
	var dp []domain.Post
	key := fmt.Sprintf("post:first:%d", id)
	val, err := p.cmd.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		p.l.Warn("缓存未命中", zap.String("key", key))
		return nil, nil
	} else if err != nil {
		p.l.Warn("缓存获取失败", zap.Error(err), zap.String("key", key))
		return nil, err
	}
	if er := json.Unmarshal(val, &dp); er != nil {
		p.l.Error("反序列化失败", zap.Error(er), zap.String("key", key))
		return nil, er
	}
	return dp, nil
}

// SetFirstPage 设置个人第一页帖子摘要
func (p *postCache) SetFirstPage(ctx context.Context, id int64, post []domain.Post) error {
	for i := 0; i < len(post); i++ {
		post[i].Content = post[i].Abstract()
	}
	val, err := json.Marshal(post)
	if err != nil {
		p.l.Error("序列化失败", zap.Error(err))
		return err
	}
	key := fmt.Sprintf("post:first:%d", id)
	return p.cmd.Set(ctx, key, val, time.Minute*15).Err()
}

// GetPubFirstPage 获取第一页公开帖子摘要
func (p *postCache) GetPubFirstPage(ctx context.Context, id int64) ([]domain.Post, error) {
	var dp []domain.Post
	key := fmt.Sprintf("post:pub:first:%d", id)
	val, err := p.cmd.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		p.l.Warn("缓存未命中", zap.String("key", key))
		return nil, nil
	} else if err != nil {
		p.l.Warn("缓存获取失败", zap.Error(err), zap.String("key", key))
		return nil, err
	}
	if er := json.Unmarshal(val, &dp); er != nil {
		p.l.Error("反序列化失败", zap.Error(er), zap.String("key", key))
		return nil, er
	}
	return dp, nil
}

// SetPubFirstPage 设置第一页公开帖子摘要
func (p *postCache) SetPubFirstPage(ctx context.Context, id int64, post []domain.Post) error {
	for i := 0; i < len(post); i++ {
		post[i].Content = post[i].Abstract()
	}
	val, err := json.Marshal(post)
	if err != nil {
		p.l.Error("序列化失败", zap.Error(err))
		return err
	}
	key := fmt.Sprintf("post:pub:first:%d", id)
	return p.cmd.Set(ctx, key, val, time.Minute*15).Err()
}

// DelFirstPage 删除第一页帖子摘要
func (p *postCache) DelFirstPage(ctx context.Context, id int64) error {
	key := fmt.Sprintf("post:first:%d", id)
	return p.cmd.Del(ctx, key).Err()
}

// DelPunFirstPage 删除第一页帖子摘要
func (p *postCache) DelPunFirstPage(ctx context.Context, id int64) error {
	key := fmt.Sprintf("post:pub:first:%d", id)
	return p.cmd.Del(ctx, key).Err()
}

// GetDetail 获取帖子详情缓存
func (p *postCache) GetDetail(ctx context.Context, id int64) (domain.Post, error) {
	var dp domain.Post
	key := fmt.Sprintf("post:detail:%d", id)
	val, err := p.cmd.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		p.l.Warn("缓存未命中", zap.String("key", key))
		return dp, nil
	} else if err != nil {
		p.l.Warn("缓存获取失败", zap.Error(err), zap.String("key", key))
		return dp, err
	}
	if er := json.Unmarshal(val, &dp); er != nil {
		p.l.Error("反序列化失败", zap.Error(er), zap.String("key", key))
		return dp, er
	}
	return dp, nil
}

// SetDetail 设置帖子详情缓存
func (p *postCache) SetDetail(ctx context.Context, post domain.Post) error {
	val, err := json.Marshal(post)
	if err != nil {
		p.l.Error("序列化失败", zap.Error(err))
		return err
	}
	key := fmt.Sprintf("post:detail:%d", post.ID)
	return p.cmd.Set(ctx, key, val, time.Minute*15).Err()
}

// GetPubDetail 获取公开帖子详情缓存
func (p *postCache) GetPubDetail(ctx context.Context, id int64) (domain.Post, error) {
	var dp domain.Post
	key := fmt.Sprintf("post:pub:detail:%d", id)
	val, err := p.cmd.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		p.l.Warn("缓存未命中", zap.String("key", key))
		return dp, nil
	} else if err != nil {
		p.l.Warn("缓存获取失败", zap.Error(err), zap.String("key", key))
		return dp, err
	}
	if er := json.Unmarshal(val, &dp); er != nil {
		p.l.Error("反序列化失败", zap.Error(er), zap.String("key", key))
		return dp, er
	}
	return dp, nil
}

// SetPubDetail 设置公开帖子详情缓存
func (p *postCache) SetPubDetail(ctx context.Context, post domain.Post) error {
	key := fmt.Sprintf("post:pub:detail:%d", post.ID)
	val, err := json.Marshal(post)
	if err != nil {
		p.l.Error("序列化失败", zap.Error(err))
		return err
	}
	return p.cmd.Set(ctx, key, val, time.Minute*15).Err()
}
