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
	Get(ctx context.Context, key string) (domain.Post, error)
	Set(ctx context.Context, key string, post domain.Post) error
	Del(ctx context.Context, key string) error
	PreHeat(ctx context.Context, posts []domain.Post) error
	GetList(ctx context.Context, page int, size int) ([]domain.Post, error)
	SetList(ctx context.Context, page int, size int, posts []domain.Post) error
	DelList(ctx context.Context, key string) error
	GetPubList(ctx context.Context, page int, size int) ([]domain.Post, error)
	SetPubList(ctx context.Context, page int, size int, posts []domain.Post) error
	DelPubList(ctx context.Context, key string) error
	GetPub(ctx context.Context, key string) (domain.Post, error)
	SetPub(ctx context.Context, key string, post domain.Post) error
	DelPub(ctx context.Context, key string) error
	SetEmpty(ctx context.Context, key string) error
	IsEmpty(ctx context.Context, key string) (bool, error)
}

type postCache struct {
	client        redis.Cmdable
	expiration    time.Duration
	prefix        string
	emptyPrefix   string
	listPrefix    string
	listPubPrefix string
	pubPrefix     string
	lockPrefix    string
}

func NewPostCache(client redis.Cmdable) PostCache {
	rand.Seed(time.Now().UnixNano()) // 初始化随机数种子
	return &postCache{
		client:        client,
		expiration:    time.Minute * 30, // 基础过期时间30分钟
		prefix:        "linkeme:post:",
		emptyPrefix:   "linkeme:post:empty:",
		listPrefix:    "linkeme:post:list:",
		listPubPrefix: "linkeme:post:list:pub:",
		pubPrefix:     "linkeme:post:pub:",
		lockPrefix:    "linkeme:post:lock:",
	}
}

// Get 获取帖子缓存
func (c *postCache) Get(ctx context.Context, key string) (domain.Post, error) {
	if key == "" {
		return domain.Post{}, fmt.Errorf("无效的帖子ID: %s", key)
	}

	// 先检查是否是空对象
	isEmpty, err := c.IsEmpty(ctx, key)
	if err == nil && isEmpty {
		return domain.Post{}, fmt.Errorf("帖子不存在 %s", key)
	}

	data, err := c.client.Get(ctx, c.key(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return domain.Post{}, fmt.Errorf("缓存中未找到帖子 %s", key)
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
func (c *postCache) Set(ctx context.Context, key string, post domain.Post) error {
	if key == "" {
		return fmt.Errorf("无效的帖子ID: %s", key)
	}

	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("序列化帖子数据失败: %v", err)
	}

	// 在基础过期时间上增加随机时间,范围是0-5分钟
	randomExpiration := c.expiration + time.Duration(rand.Int63n(300))*time.Second
	err = c.client.Set(ctx, c.key(key), data, randomExpiration).Err()
	if err != nil {
		return fmt.Errorf("设置缓存失败: %v", err)
	}

	// 删除空对象标记(如果存在)
	_ = c.client.Del(ctx, c.emptyKey(key))

	return nil
}

// SetEmpty 缓存空对象,使用较短的过期时间
func (c *postCache) SetEmpty(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("无效的帖子ID: %s", key)
	}

	// 空对象使用空字符串标记，不需要存储完整结构
	return c.client.Set(ctx, c.emptyKey(key), "1", time.Minute*5).Err()
}

// IsEmpty 检查是否存在空对象标记
func (c *postCache) IsEmpty(ctx context.Context, key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("无效的帖子ID: %s", key)
	}

	exists, err := c.client.Exists(ctx, c.emptyKey(key)).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

// PreHeat 预热热点缓存,设置更长的过期时间
func (c *postCache) PreHeat(ctx context.Context, posts []domain.Post) error {
	if len(posts) == 0 {
		return fmt.Errorf("帖子列表为空")
	}

	// 使用Pipeline批量设置
	pipe := c.client.Pipeline()

	// 预热单个帖子缓存
	for _, post := range posts {
		if post.ID == 0 {
			continue
		}

		data, err := json.Marshal(post)
		if err != nil {
			continue
		}

		// 热点数据设置更长过期时间
		pipe.Set(ctx, c.pubKey(fmt.Sprint(post.ID)), data, time.Hour*2)
	}

	// 预热列表缓存
	postsData, err := json.Marshal(posts)
	if err == nil {
		// 只缓存前三页,每页设置不同的过期时间
		for i := 1; i <= 3; i++ {
			key := c.pubListKey(fmt.Sprintf("%d_%d", i, 10)) // 默认每页10条
			expiration := time.Hour * time.Duration(4-i)     // 第一页3小时,第二页2小时,第三页1小时
			pipe.Set(ctx, key, postsData, expiration)
		}
	}

	// 执行批量操作
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("预热缓存失败: %v", err)
	}

	return nil
}

// Del 删除帖子缓存
func (c *postCache) Del(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("无效的帖子ID: %s", key)
	}

	pipe := c.client.Pipeline()
	pipe.Del(ctx, c.key(key))
	pipe.Del(ctx, c.emptyKey(key))
	pipe.Del(ctx, c.pubKey(key))
	_, err := pipe.Exec(ctx)
	return err
}

// GetList 获取帖子列表缓存
func (c *postCache) GetList(ctx context.Context, page int, size int) ([]domain.Post, error) {
	if page <= 0 || size <= 0 {
		return nil, fmt.Errorf("无效的分页参数: page=%d, size=%d", page, size)
	}

	key := fmt.Sprintf("%d_%d", page, size)
	data, err := c.client.Get(ctx, c.listKey(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("缓存中未找到帖子列表: %s", key)
		}
		return nil, fmt.Errorf("从缓存获取帖子列表失败: %v", err)
	}

	var posts []domain.Post
	if err = json.Unmarshal(data, &posts); err != nil {
		// 删除可能损坏的缓存数据
		_ = c.client.Del(ctx, c.listKey(key))
		return nil, fmt.Errorf("反序列化帖子列表数据失败: %v", err)
	}

	return posts, nil
}

// SetList 设置帖子列表缓存
func (c *postCache) SetList(ctx context.Context, page int, size int, posts []domain.Post) error {
	if page <= 0 || size <= 0 {
		return fmt.Errorf("无效的分页参数: page=%d, size=%d", page, size)
	}

	data, err := json.Marshal(posts)
	if err != nil {
		return fmt.Errorf("序列化帖子列表数据失败: %v", err)
	}

	// 列表缓存使用随机过期时间
	randomExpiration := c.expiration + time.Duration(rand.Int63n(300))*time.Second
	key := fmt.Sprintf("%d_%d", page, size)
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
func (c *postCache) GetPubList(ctx context.Context, page int, size int) ([]domain.Post, error) {
	if page <= 0 || size <= 0 {
		return nil, fmt.Errorf("无效的分页参数: page=%d, size=%d", page, size)
	}

	key := fmt.Sprintf("%d_%d", page, size)
	data, err := c.client.Get(ctx, c.pubListKey(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("缓存中未找到已发布帖子列表: %s", key)
		}
		return nil, fmt.Errorf("从缓存获取已发布帖子列表失败: %v", err)
	}

	var posts []domain.Post
	if err = json.Unmarshal(data, &posts); err != nil {
		// 删除可能损坏的缓存数据
		_ = c.client.Del(ctx, c.pubListKey(key))
		return nil, fmt.Errorf("反序列化已发布帖子列表数据失败: %v", err)
	}

	return posts, nil
}

// SetPubList 设置已发布帖子列表缓存
func (c *postCache) SetPubList(ctx context.Context, page int, size int, posts []domain.Post) error {
	if page <= 0 || size <= 0 {
		return fmt.Errorf("无效的分页参数: page=%d, size=%d", page, size)
	}

	// 获取分布式锁
	key := fmt.Sprintf("%d_%d", page, size)
	lockKey := c.lockKey("publist:" + key)
	success, err := c.acquireLock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return fmt.Errorf("获取分布式锁失败: %v", err)
	}
	if !success {
		return fmt.Errorf("获取分布式锁失败: 锁已被占用")
	}
	defer c.releaseLock(ctx, lockKey)

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
		var cursor uint64
		var keys []string

		for {
			var scanKeys []string
			var err error

			// 每次扫描100个key
			scanKeys, cursor, err = c.client.Scan(ctx, cursor, pattern, 100).Result()
			if err != nil {
				return fmt.Errorf("扫描缓存键失败: %v", err)
			}

			keys = append(keys, scanKeys...)

			// 如果游标为0，表示扫描完成
			if cursor == 0 {
				break
			}
		}

		// 如果有key需要删除,则使用pipeline批量删除
		if len(keys) > 0 {
			var err error
			for i := 0; i < 3; i++ { // 重试3次
				pipe := c.client.Pipeline()
				for _, key := range keys {
					pipe.Del(ctx, key)
				}

				if _, err = pipe.Exec(ctx); err == nil {
					break
				}

				if i == 2 { // 最后一次重试失败
					return fmt.Errorf("批量删除缓存键失败(重试3次): %v", err)
				}
				time.Sleep(time.Millisecond * 50) // 失败后短暂等待再重试
			}
		}
		return nil
	}
	return c.client.Del(ctx, c.pubListKey(key)).Err()
}

// GetPub 获取已发布帖子缓存
func (c *postCache) GetPub(ctx context.Context, key string) (domain.Post, error) {
	if key == "" {
		return domain.Post{}, fmt.Errorf("无效的帖子ID: %s", key)
	}

	// 先检查是否是空对象
	isEmpty, err := c.IsEmpty(ctx, key)
	if err == nil && isEmpty {
		return domain.Post{}, fmt.Errorf("帖子不存在 %s", key)
	}

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
	if key == "" {
		return fmt.Errorf("无效的帖子ID: %s", key)
	}

	// 获取分布式锁
	lockKey := c.lockKey("pub:" + key)
	success, err := c.acquireLock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return fmt.Errorf("获取分布式锁失败: %v", err)
	}
	if !success {
		return fmt.Errorf("获取分布式锁失败: 锁已被占用")
	}
	defer c.releaseLock(ctx, lockKey)

	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("序列化已发布帖子数据失败: %v", err)
	}

	// 使用随机过期时间
	randomExpiration := c.expiration + time.Duration(rand.Int63n(300))*time.Second
	err = c.client.Set(ctx, c.pubKey(key), data, randomExpiration).Err()
	if err != nil {
		return fmt.Errorf("设置缓存失败: %v", err)
	}

	// 删除空对象标记(如果存在)
	_ = c.client.Del(ctx, c.emptyKey(key))

	return nil
}

// DelPub 删除已发布帖子缓存
func (c *postCache) DelPub(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("无效的帖子ID: %s", key)
	}

	return c.client.Del(ctx, c.pubKey(key)).Err()
}

// 辅助方法: 获取分布式锁
func (c *postCache) acquireLock(ctx context.Context, lockKey string, expiration time.Duration) (bool, error) {
	return c.client.SetNX(ctx, lockKey, "1", expiration).Result()
}

// 辅助方法: 释放分布式锁
func (c *postCache) releaseLock(ctx context.Context, lockKey string) {
	_ = c.client.Del(ctx, lockKey)
}

func (c *postCache) key(key string) string {
	return c.prefix + key
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

func (c *postCache) emptyKey(key string) string {
	return c.emptyPrefix + key
}

func (c *postCache) lockKey(key string) string {
	return c.lockPrefix + key
}
