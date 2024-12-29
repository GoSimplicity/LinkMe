package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/redis/go-redis/v9"
)

type PostCache interface {
	Get(ctx context.Context, id int64) (domain.Post, error)
	Set(ctx context.Context, post domain.Post) error
	Del(ctx context.Context, id int64) error
	PreHeat(ctx context.Context, post domain.Post) error
	GetList(ctx context.Context, key string) ([]domain.Post, error)
	SetList(ctx context.Context, key string, posts []domain.Post) error
	DelList(ctx context.Context, key string) error
	GetPubList(ctx context.Context, key string) ([]domain.Post, error)
	SetPubList(ctx context.Context, key string, posts []domain.Post) error
	DelPubList(ctx context.Context, key string) error
	GetPub(ctx context.Context, key string) (domain.Post, error)
	SetPub(ctx context.Context, key string, post domain.Post) error
	DelPub(ctx context.Context, key string) error
}

type postCache struct {
	client        redis.Cmdable
	expiration    time.Duration
	prefix        string
	listPrefix    string
	listPubPrefix string
	pubPrefix     string
}

func NewPostCache(client redis.Cmdable) PostCache {
	return &postCache{
		client:        client,
		expiration:    time.Minute * 30, // 基础过期时间30分钟
		prefix:        "linkeme:post:",
		listPrefix:    "linkeme:post:list:",
		listPubPrefix: "linkeme:post:list:pub:",
		pubPrefix:     "linkeme:post:pub:",
	}
}

// Get 获取帖子缓存
func (c *postCache) Get(ctx context.Context, id int64) (domain.Post, error) {
	if id <= 0 {
		return domain.Post{}, fmt.Errorf("无效的帖子ID: %d", id)
	}

	key := c.key(id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return domain.Post{}, fmt.Errorf("缓存中未找到帖子 %d", id)
		}
		return domain.Post{}, fmt.Errorf("从缓存获取帖子失败: %v", err)
	}

	var post domain.Post
	if err = json.Unmarshal(data, &post); err != nil {
		return domain.Post{}, fmt.Errorf("反序列化帖子数据失败: %v", err)
	}

	return post, nil
}

// Set 设置帖子缓存,使用随机过期时间防止缓存雪崩
func (c *postCache) Set(ctx context.Context, post domain.Post) error {
	if post.ID <= 0 {
		return fmt.Errorf("无效的帖子ID: %d", post.ID)
	}

	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("序列化帖子数据失败: %v", err)
	}

	// 在基础过期时间上增加随机时间,范围是0-5分钟
	randomExpiration := c.expiration + time.Duration(rand.Int63n(300))*time.Second
	key := c.key(int64(post.ID))
	return c.client.Set(ctx, key, data, randomExpiration).Err()
}

// PreHeat 预热热点缓存,设置更长的过期时间
func (c *postCache) PreHeat(ctx context.Context, post domain.Post) error {
	if post.ID <= 0 {
		return fmt.Errorf("无效的帖子ID: %d", post.ID)
	}

	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("序列化帖子数据失败: %v", err)
	}

	// 热点数据设置更长的过期时间,1小时
	key := c.key(int64(post.ID))
	return c.client.Set(ctx, key, data, time.Hour).Err()
}

// Del 删除帖子缓存
func (c *postCache) Del(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("无效的帖子ID: %d", id)
	}

	key := c.key(id)
	return c.client.Del(ctx, key).Err()
}

// GetList 获取帖子列表缓存
func (c *postCache) GetList(ctx context.Context, key string) ([]domain.Post, error) {
	data, err := c.client.Get(ctx, c.listKey(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("缓存中未找到帖子列表: %s", key)
		}
		return nil, fmt.Errorf("从缓存获取帖子列表失败: %v", err)
	}

	var posts []domain.Post
	if err = json.Unmarshal(data, &posts); err != nil {
		return nil, fmt.Errorf("反序列化帖子列表数据失败: %v", err)
	}

	return posts, nil
}

// SetList 设置帖子列表缓存
func (c *postCache) SetList(ctx context.Context, key string, posts []domain.Post) error {
	data, err := json.Marshal(posts)
	if err != nil {
		return fmt.Errorf("序列化帖子列表数据失败: %v", err)
	}

	// 列表缓存使用随机过期时间
	randomExpiration := c.expiration + time.Duration(rand.Int63n(300))*time.Second
	return c.client.Set(ctx, c.listKey(key), data, randomExpiration).Err()
}

// DelList 删除帖子列表缓存
func (c *postCache) DelList(ctx context.Context, key string) error {
	if key == "*" {
		pattern := c.listPrefix + key
		// 使用SCAN命令批量删除,每次扫描100个key
		var cursor uint64
		var keys []string

		for {
			var scanKeys []string
			var err error
			scanKeys, cursor, err = c.client.Scan(ctx, cursor, pattern, 100).Result()
			if err != nil {
				return fmt.Errorf("扫描缓存键失败: %v", err)
			}

			keys = append(keys, scanKeys...)

			if cursor == 0 {
				break
			}
		}

		// 如果有key需要删除,则使用pipeline批量删除
		if len(keys) > 0 {
			pipe := c.client.Pipeline()
			for _, key := range keys {
				pipe.Del(ctx, key)
			}
			if _, err := pipe.Exec(ctx); err != nil {
				return fmt.Errorf("批量删除缓存键失败: %v", err)
			}
		}

		return nil
	}
	return c.client.Del(ctx, c.listKey(key)).Err()
}

// GetPubList 获取已发布帖子列表缓存
func (c *postCache) GetPubList(ctx context.Context, key string) ([]domain.Post, error) {
	data, err := c.client.Get(ctx, c.pubListKey(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("缓存中未找到已发布帖子列表: %s", key)
		}
		return nil, fmt.Errorf("从缓存获取已发布帖子列表失败: %v", err)
	}

	var posts []domain.Post
	if err = json.Unmarshal(data, &posts); err != nil {
		return nil, fmt.Errorf("反序列化已发布帖子列表数据失败: %v", err)
	}

	return posts, nil
}

// SetPubList 设置已发布帖子列表缓存
func (c *postCache) SetPubList(ctx context.Context, key string, posts []domain.Post) error {
	data, err := json.Marshal(posts)
	if err != nil {
		return fmt.Errorf("序列化已发布帖子列表数据失败: %v", err)
	}

	// 列表缓存使用随机过期时间
	randomExpiration := c.expiration + time.Duration(rand.Int63n(300))*time.Second
	return c.client.Set(ctx, c.pubListKey(key), data, randomExpiration).Err()
}

// DelPubList 删除已发布帖子列表缓存
func (c *postCache) DelPubList(ctx context.Context, key string) error {
	if key == "*" {
		pattern := c.listPubPrefix + key
		// 使用SCAN命令批量删除,每次扫描100个key
		var cursor uint64
		var keys []string

		for {
			var scanKeys []string
			var err error
			scanKeys, cursor, err = c.client.Scan(ctx, cursor, pattern, 100).Result()
			if err != nil {
				return fmt.Errorf("扫描缓存键失败: %v", err)
			}

			keys = append(keys, scanKeys...)

			if cursor == 0 {
				break
			}
		}

		// 如果有key需要删除,则使用pipeline批量删除
		if len(keys) > 0 {
			pipe := c.client.Pipeline()
			for _, key := range keys {
				pipe.Del(ctx, key)
			}
			if _, err := pipe.Exec(ctx); err != nil {
				return fmt.Errorf("批量删除缓存键失败: %v", err)
			}
		}

		return nil
	}
	return c.client.Del(ctx, c.pubListKey(key)).Err()
}

// GetPub 获取已发布帖子缓存
func (c *postCache) GetPub(ctx context.Context, key string) (domain.Post, error) {
	data, err := c.client.Get(ctx, c.pubKey(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return domain.Post{}, fmt.Errorf("缓存中未找到已发布帖子: %s", key)
		}
		return domain.Post{}, fmt.Errorf("从缓存获取已发布帖子失败: %v", err)
	}

	var post domain.Post
	if err = json.Unmarshal(data, &post); err != nil {
		return domain.Post{}, fmt.Errorf("反序列化已发布帖子数据失败: %v", err)
	}

	return post, nil
}

// SetPub 设置已发布帖子缓存
func (c *postCache) SetPub(ctx context.Context, key string, post domain.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("序列化已发布帖子数据失败: %v", err)
	}

	// 使用随机过期时间
	randomExpiration := c.expiration + time.Duration(rand.Int63n(300))*time.Second
	return c.client.Set(ctx, c.pubKey(key), data, randomExpiration).Err()
}

// DelPub 删除已发布帖子缓存
func (c *postCache) DelPub(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.pubKey(key)).Err()
}

func (c *postCache) key(id int64) string {
	return c.prefix + fmt.Sprint(id)
}

func (c *postCache) listKey(key string) string {
	return c.listPrefix + key
}

func (c *postCache) pubListKey(key string) string {
	return c.listPubPrefix + key
}

func (c *postCache) pubKey(key string) string {
	return c.pubPrefix + key
}
