package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/redis/go-redis/v9"
)

type CommentCache interface {
	Get(ctx context.Context, postId int64) (domain.Comment, error)
	Set(ctx context.Context, du domain.Comment) error
}

type commentCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewCommentCache(cmd redis.Cmdable) CommentCache {
	return &commentCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

// Get 从redis中获取数据并反序列化
func (u *commentCache) Get(ctx context.Context, postId int64) (domain.Comment, error) {
	var dc domain.Comment
	key := fmt.Sprintf("linkme:comment:%d", postId)
	// 从redis中读取数据
	data, err := u.cmd.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.Comment{}, err
		}
		return domain.Comment{}, err
	}
	if err = json.Unmarshal([]byte(data), &dc); err != nil {
		return domain.Comment{}, fmt.Errorf("反序列化评论数据失败: %v", err)
	}
	return dc, nil
}

// Set 将传入的du结构体序列化存入redis中
func (u *commentCache) Set(ctx context.Context, dc domain.Comment) error {
	key := fmt.Sprintf("linkme:comment:%d", dc.PostId)
	data, err := json.Marshal(dc)
	if err != nil {
		return fmt.Errorf("序列化评论数据失败: %v", err)
	}

	if err := u.cmd.Set(ctx, key, data, u.expiration).Err(); err != nil {
		return fmt.Errorf("缓存评论数据失败: %v", err)
	}
	return nil
}
