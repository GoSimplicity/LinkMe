/*
 * MIT License
 *
 * Copyright (c) 2024 Bamboo
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CoreCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	MGet(ctx context.Context, keys ...string) (map[string][]byte, error)
	MSet(ctx context.Context, items map[string][]byte, expiration time.Duration) error
	MDelete(ctx context.Context, keys ...string) error

	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)

	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)

	Lock(ctx context.Context, key string, expiration time.Duration) (bool, error)
	Unlock(ctx context.Context, key string) error

	Ping(ctx context.Context) error

	Close() error
}

type coreCache struct {
	redis redis.Cmdable
}

func NewCoreCache(redis redis.Cmdable) CoreCache {
	return &coreCache{
		redis: redis,
	}
}

// Close 关闭redis客户端
func (c *coreCache) Close() error {
	if client, ok := c.redis.(*redis.Client); ok {
		return client.Close()
	}
	return nil
}

// Decr 递减
func (c *coreCache) Decr(ctx context.Context, key string) (int64, error) {
	return c.redis.Decr(ctx, key).Result()
}

// Delete 删除缓存
func (c *coreCache) Delete(ctx context.Context, key string) error {
	return c.redis.Del(ctx, key).Err()
}

// Eval 执行lua脚本
func (c *coreCache) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return c.redis.Eval(ctx, script, keys, args...).Result()
}

// Exists 判断key是否存在
func (c *coreCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.redis.Exists(ctx, key).Result()
	return result > 0, err
}

// Get 获取缓存
func (c *coreCache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.redis.Get(ctx, key).Bytes()
}

// Incr 递增
func (c *coreCache) Incr(ctx context.Context, key string) (int64, error) {
	return c.redis.Incr(ctx, key).Result()
}

// Lock 设置分布式锁
func (c *coreCache) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	// 使用SET NX实现分布式锁
	lockKey := "lock:" + key
	// 生成随机值作为锁标识
	lockValue := time.Now().String()
	result, err := c.redis.SetNX(ctx, lockKey, lockValue, expiration).Result()
	return result, err
}

// MDelete 批量删除缓存
func (c *coreCache) MDelete(ctx context.Context, keys ...string) error {
	return c.redis.Del(ctx, keys...).Err()
}

// MGet 批量获取缓存
func (c *coreCache) MGet(ctx context.Context, keys ...string) (map[string][]byte, error) {
	values, err := c.redis.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]byte)
	for idx, key := range keys {
		if values[idx] != nil {
			if str, ok := values[idx].(string); ok {
				result[key] = []byte(str)
			}
		}
	}
	return result, nil
}

// MSet 批量设置缓存
func (c *coreCache) MSet(ctx context.Context, items map[string][]byte, expiration time.Duration) error {
	if len(items) == 0 {
		return nil
	}

	// 转换为redis.MSet所需的参数格式
	pairs := make([]interface{}, 0, len(items)*2)
	for k, v := range items {
		pairs = append(pairs, k, v)
	}

	pipe := c.redis.Pipeline()
	pipe.MSet(ctx, pairs...)

	// 设置过期时间
	if expiration > 0 {
		for k := range items {
			pipe.Expire(ctx, k, expiration)
		}
	}

	_, err := pipe.Exec(ctx)
	return err
}

// Ping 健康检查
func (c *coreCache) Ping(ctx context.Context) error {
	return c.redis.Ping(ctx).Err()
}

// Set 设置缓存
func (c *coreCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return c.redis.Set(ctx, key, value, expiration).Err()
}

// Unlock 释放分布式锁
func (c *coreCache) Unlock(ctx context.Context, key string) error {
	lockKey := "lock:" + key
	return c.redis.Del(ctx, lockKey).Err()
}
