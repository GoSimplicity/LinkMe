package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, du domain.User) error
}

type userCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &userCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

// Get 从redis中获取数据并反序列化
func (u *userCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	if uid <= 0 {
		return domain.User{}, fmt.Errorf("invalid user id: %d", uid)
	}

	var du domain.User
	key := fmt.Sprintf("linkeme:user:%d", uid)

	// 从redis中读取数据
	data, err := u.cmd.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.User{}, fmt.Errorf("user %d not found in cache", uid)
		}
		return domain.User{}, fmt.Errorf("failed to get user from cache: %v", err)
	}

	if err = json.Unmarshal([]byte(data), &du); err != nil {
		return domain.User{}, fmt.Errorf("failed to unmarshal user data: %v", err)
	}

	// 如果用户已被删除,则不返回数据
	if du.Deleted {
		return domain.User{}, fmt.Errorf("user %d has been deleted", uid)
	}

	return du, nil
}

// Set 将传入的du结构体序列化存入redis中
func (u *userCache) Set(ctx context.Context, du domain.User) error {
	if du.ID <= 0 {
		return fmt.Errorf("invalid user id: %d", du.ID)
	}

	key := fmt.Sprintf("linkme:user:%d", du.ID)
	data, err := json.Marshal(du)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %v", err)
	}

	// 向redis中插入数据
	if err = u.cmd.Set(ctx, key, data, u.expiration).Err(); err != nil {
		return fmt.Errorf("failed to set user in cache: %v", err)
	}

	return nil
}
