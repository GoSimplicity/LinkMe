package cache

import (
	"LinkMe/internal/domain"
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type PostCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Post, error)   // 获取用户的第一页帖子缓存
	SetFirstPage(ctx context.Context, uid int64, res []domain.Post) error // 设置用户的第一页帖子缓存
	DelFirstPage(ctx context.Context, uid int64) error                    // 删除用户的第一页帖子缓存
	Get(ctx context.Context, id int64) (domain.Post, error)               // 根据ID获取一个帖子缓存
	Set(ctx context.Context, art domain.Post) error                       // 设置一个帖子缓存
	GetPub(ctx context.Context, id int64) (domain.Post, error)            // 根据ID获取一个已发布的帖子缓存
	SetPub(ctx context.Context, res domain.Post) error                    // 设置一个已发布的帖子缓存
}

type postCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
	l          *zap.Logger
}

func NewPostCache(cmd redis.Cmdable, l *zap.Logger) PostCache {
	return &postCache{
		cmd:        cmd,
		expiration: time.Minute * 10,
		l:          l,
	}
}
func (p *postCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postCache) SetFirstPage(ctx context.Context, uid int64, res []domain.Post) error {
	//TODO implement me
	panic("implement me")
}

func (p *postCache) DelFirstPage(ctx context.Context, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (p *postCache) Get(ctx context.Context, id int64) (domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postCache) Set(ctx context.Context, art domain.Post) error {
	//TODO implement me
	panic("implement me")
}

func (p *postCache) GetPub(ctx context.Context, id int64) (domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postCache) SetPub(ctx context.Context, res domain.Post) error {
	//TODO implement me
	panic("implement me")
}
