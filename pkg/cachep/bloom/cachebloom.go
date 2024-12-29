package bloom

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/redis/go-redis/v9"
)

type CacheBloom struct {
	bf     *bloom.BloomFilter
	client redis.Cmdable
}

func NewCacheBloom(client redis.Cmdable) *CacheBloom {
	return &CacheBloom{
		bf:     bloom.NewWithEstimates(1000000, 0.01),
		client: client,
	}
}

// QueryData 查询数据，如果缓存和布隆过滤器都没有命中，则存储并返回传递的数据
func QueryData[T any](cb *CacheBloom, ctx context.Context, key string, data T, ttl time.Duration) (T, error) {
	if cb == nil {
		var zeroValue T
		return zeroValue, errors.New("CacheBloom实例不能为空")
	}

	if key == "" {
		var zeroValue T
		return zeroValue, errors.New("键不能为空")
	}

	// 检查布隆过滤器
	if !cb.bf.TestString(key) {
		if !isNil(data) {
			return CacheData(cb, ctx, key, data, ttl)
		}
		var zeroValue T
		return zeroValue, nil
	}

	// 检查缓存
	cachedData, err := cb.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			if !isNil(data) {
				return CacheData(cb, ctx, key, data, ttl)
			}
			var zeroValue T
			return zeroValue, nil
		}
		var zeroValue T
		return zeroValue, err
	}

	var result T
	if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
		var zeroValue T
		return zeroValue, err
	}

	return result, nil
}

// CacheData 序列化数据并存储到 Redis 和布隆过滤器中
func CacheData[T any](cb *CacheBloom, ctx context.Context, key string, data T, ttl time.Duration) (T, error) {
	if key == "" {
		return data, errors.New("键不能为空")
	}

	if ttl <= 0 {
		return data, errors.New("过期时间必须为正数")
	}

	serializedData, err := json.Marshal(data)
	if err != nil {
		return data, err
	}

	// 使用pipeline优化Redis操作
	pipe := cb.client.Pipeline()
	pipe.Set(ctx, key, serializedData, ttl)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return data, err
	}

	cb.bf.AddString(key)
	return data, nil
}

// SetEmptyCache 缓存空对象，防止缓存穿透
func (cb *CacheBloom) SetEmptyCache(ctx context.Context, key string, ttl time.Duration) error {
	if key == "" {
		return errors.New("键不能为空")
	}

	if ttl <= 0 {
		return errors.New("过期时间必须为正数")
	}

	emptyValue, err := json.Marshal(struct{}{})
	if err != nil {
		return err
	}

	pipe := cb.client.Pipeline()
	pipe.Set(ctx, key, emptyValue, ttl)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	cb.bf.AddString(key)
	return nil
}

// RebuildBloomFilter 重建布隆过滤器
// 导出该方法以便外部调用
func (cb *CacheBloom) RebuildBloomFilter(ctx context.Context) error {
	if cb == nil {
		return errors.New("CacheBloom实例不能为空")
	}

	newBF := bloom.NewWithEstimates(1000000, 0.01)

	// 使用SCAN命令替代KEYS,避免阻塞Redis
	iter := cb.client.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		newBF.AddString(iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}

	cb.bf = newBF
	return nil
}

// isNil 使用反射检查值是否为 nil
func isNil(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface())
	}
}
