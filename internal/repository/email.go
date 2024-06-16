package repository

import (
	"LinkMe/internal/repository/cache"
	"context"
)

type EmailRepository interface {
	CheckCode(ctx context.Context, email, vCode string) (bool, error)
	StoreVCode(ctx context.Context, email, vCode string) error
}

// emailRepository 实现了 EmailRepository 接口
type emailRepository struct {
	cache cache.EmailCache
}

// NewEmailRepository 创建并返回一个新的 smsRepository 实例
func NewEmailRepository(cache cache.EmailCache) EmailRepository {
	return &emailRepository{
		cache: cache,
	}
}

func (e emailRepository) CheckCode(ctx context.Context, email, vCode string) (bool, error) {
	storedCode, err := e.cache.GetVCode(ctx, email)
	return storedCode == vCode, err
}

func (e emailRepository) StoreVCode(ctx context.Context, email, vCode string) error {
	return e.cache.StoreVCode(ctx, email, vCode)
}
