package cache

import (
	"LinkMe/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
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
		expiration: time.Minute * 10,
	}
}

// Get 从redis中获取数据并反序列化
func (u *userCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	var du domain.User
	key := fmt.Sprintf("linkeme:user:%d", uid)
	// 从redis中读取数据
	data, err := u.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	err = json.Unmarshal([]byte(data), &du)
	return du, err
}

// Set 将传入的du结构体序列化存入redis中
func (u *userCache) Set(ctx context.Context, du domain.User) error {
	key := fmt.Sprintf("linkme:user:%d", du.ID)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	// 向redis中插入数据
	return u.cmd.Set(ctx, key, data, u.expiration).Err()
}
