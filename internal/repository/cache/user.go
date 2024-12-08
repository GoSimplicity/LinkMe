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
		return domain.User{}, fmt.Errorf("无效的用户ID: %d", uid)
	}

	var du domain.User
	key := fmt.Sprintf("linkeme:user:%d", uid)

	// 从redis中读取数据
	data, err := u.cmd.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.User{}, fmt.Errorf("缓存中未找到用户 %d", uid)
		}
		return domain.User{}, fmt.Errorf("从缓存获取用户失败: %v", err)
	}

	if err = json.Unmarshal([]byte(data), &du); err != nil {
		return domain.User{}, fmt.Errorf("反序列化用户数据失败: %v", err)
	}

	// 如果用户已被删除,则不返回数据
	if du.Deleted {
		return domain.User{}, fmt.Errorf("用户 %d 已被删除", uid)
	}

	return du, nil
}

// Set 将传入的du结构体序列化存入redis中
func (u *userCache) Set(ctx context.Context, du domain.User) error {
	if du.ID <= 0 {
		return fmt.Errorf("无效的用户ID: %d", du.ID)
	}

	key := fmt.Sprintf("linkme:user:%d", du.ID)
	data, err := json.Marshal(du)
	if err != nil {
		return fmt.Errorf("序列化用户数据失败: %v", err)
	}

	// 向redis中插入数据,使用pipeline优化性能
	pipe := u.cmd.Pipeline()
	pipe.Set(ctx, key, data, u.expiration)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("缓存用户数据失败: %v", err)
	}

	return nil
}
