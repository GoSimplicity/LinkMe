package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type LotteryDrawCache interface {
	// GetLotteryDraw 从缓存中获取指定ID的抽奖活动
	GetLotteryDraw(ctx context.Context, id int) (domain.LotteryDraw, error)
	// SetLotteryDraw 将抽奖活动设置到缓存中
	SetLotteryDraw(ctx context.Context, draw domain.LotteryDraw) error
	// DeleteLotteryDraw 从缓存中删除指定ID的抽奖活动
	DeleteLotteryDraw(ctx context.Context, id int) error
	// GetLotteryDrawWithLock 在缓存未命中时，使用分布式锁从数据库获取并设置缓存
	GetLotteryDrawWithLock(ctx context.Context, id int, fetchFromDB func() (domain.LotteryDraw, error)) (domain.LotteryDraw, error)

	// GetSecondKillEvent 从缓存中获取指定ID的秒杀活动
	GetSecondKillEvent(ctx context.Context, id int) (domain.SecondKillEvent, error)
	// SetSecondKillEvent 将秒杀活动设置到缓存中
	SetSecondKillEvent(ctx context.Context, event domain.SecondKillEvent) error
	// DeleteSecondKillEvent 从缓存中删除指定ID的秒杀活动
	DeleteSecondKillEvent(ctx context.Context, id int) error
	// GetSecondKillEventWithLock 在缓存未命中时，使用分布式锁从数据库获取并设置缓存
	GetSecondKillEventWithLock(ctx context.Context, id int, fetchFromDB func() (domain.SecondKillEvent, error)) (domain.SecondKillEvent, error)
}

type lotteryDrawCache struct {
	client redis.Cmdable
	logger *zap.Logger
}

func NewLotteryDrawCache(client redis.Cmdable, logger *zap.Logger) LotteryDrawCache {
	return &lotteryDrawCache{
		client: client,
		logger: logger,
	}
}

// 缓存键的前缀
const (
	lotteryDrawKeyPrefix     = "linkme:lottery_draw:"
	secondKillEventKeyPrefix = "linkme:second_kill_event:"
	lockKeyPrefix            = "lock:"
)

// 缓存数据的过期时间
const (
	lotteryDrawTTL     = 10 * time.Minute
	secondKillEventTTL = 10 * time.Minute
	lockTTL            = 5 * time.Second // 锁的过期时间，防止死锁
	maxLockRetries     = 3
	lockRetryDelay     = 100 * time.Millisecond
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

// GetLotteryDraw 从缓存中获取指定ID的抽奖活动
func (c *lotteryDrawCache) GetLotteryDraw(ctx context.Context, id int) (domain.LotteryDraw, error) {
	key := generateLotteryDrawKey(id)

	// 获取缓存数据
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存未命中
			return domain.LotteryDraw{}, ErrCacheMiss
		}
		// 其他错误
		c.logger.Error("获取抽奖活动缓存失败", zap.Error(err), zap.Int("ID", id))
		return domain.LotteryDraw{}, err
	}

	// 反序列化 JSON 数据
	var draw domain.LotteryDraw
	if err := json.Unmarshal([]byte(data), &draw); err != nil {
		c.logger.Error("反序列化抽奖活动数据失败", zap.Error(err), zap.String("Data", data))
		return domain.LotteryDraw{}, err
	}

	return draw, nil
}

// SetLotteryDraw 将抽奖活动设置到缓存中
func (c *lotteryDrawCache) SetLotteryDraw(ctx context.Context, draw domain.LotteryDraw) error {
	key := generateLotteryDrawKey(draw.ID)

	// 序列化为 JSON
	data, err := json.Marshal(draw)
	if err != nil {
		c.logger.Error("序列化抽奖活动数据失败", zap.Error(err), zap.Any("Draw", draw))
		return err
	}

	// 设置缓存数据并设置过期时间
	if err := c.client.Set(ctx, key, data, lotteryDrawTTL).Err(); err != nil {
		c.logger.Error("设置抽奖活动缓存失败", zap.Error(err), zap.Int("ID", draw.ID))
		return err
	}

	return nil
}

// DeleteLotteryDraw 从缓存中删除指定ID的抽奖活动
func (c *lotteryDrawCache) DeleteLotteryDraw(ctx context.Context, id int) error {
	key := generateLotteryDrawKey(id)

	// 删除缓存键
	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.Error("删除抽奖活动缓存失败", zap.Error(err), zap.Int("ID", id))
		return err
	}

	return nil
}

// GetLotteryDrawWithLock 在缓存未命中时，使用分布式锁从数据库获取并设置缓存
func (c *lotteryDrawCache) GetLotteryDrawWithLock(ctx context.Context, id int, fetchFromDB func() (domain.LotteryDraw, error)) (domain.LotteryDraw, error) {
	// 尝试从缓存中获取数据
	draw, err := c.GetLotteryDraw(ctx, id)
	if err == nil {
		return draw, nil
	}
	if !errors.Is(err, ErrCacheMiss) {
		return domain.LotteryDraw{}, err
	}

	// 缓存未命中，尝试获取分布式锁
	lockKey := generateLockKey(generateLotteryDrawKey(id))
	lockValue := uuid.New().String()

	// 尝试获取锁
	acquired, err := c.acquireLock(ctx, lockKey, lockValue, lockTTL)
	if err != nil {
		c.logger.Error("获取分布式锁失败", zap.Error(err), zap.String("LockKey", lockKey))
		return domain.LotteryDraw{}, err
	}

	if !acquired {
		// 获取锁失败，可能其他请求正在设置缓存，等待并重试获取缓存
		time.Sleep(lockRetryDelay)
		return c.GetLotteryDraw(ctx, id)
	}

	// 确保在函数退出时释放锁
	defer func() {
		if err := c.releaseLock(ctx, lockKey, lockValue); err != nil {
			c.logger.Error("释放分布式锁失败", zap.Error(err), zap.String("LockKey", lockKey))
		}
	}()

	// 再次尝试从缓存中获取数据，以防止锁获取后，其他请求已经设置了缓存
	draw, err = c.GetLotteryDraw(ctx, id)
	if err == nil {
		return draw, nil
	}
	if !errors.Is(err, ErrCacheMiss) {
		return domain.LotteryDraw{}, err
	}

	// 从数据库获取数据
	draw, err = fetchFromDB()
	if err != nil {
		return domain.LotteryDraw{}, err
	}

	// 设置缓存
	if err := c.SetLotteryDraw(ctx, draw); err != nil {
		c.logger.Error("设置抽奖活动缓存失败", zap.Error(err), zap.Int("ID", id))
	}

	return draw, nil
}

// GetSecondKillEvent 从缓存中获取指定ID的秒杀活动
func (c *lotteryDrawCache) GetSecondKillEvent(ctx context.Context, id int) (domain.SecondKillEvent, error) {
	key := generateSecondKillEventKey(id)

	// 获取缓存数据
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存未命中
			return domain.SecondKillEvent{}, ErrCacheMiss
		}
		// 其他错误
		c.logger.Error("获取秒杀活动缓存失败", zap.Error(err), zap.Int("ID", id))
		return domain.SecondKillEvent{}, err
	}

	// 反序列化 JSON 数据
	var event domain.SecondKillEvent
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		c.logger.Error("反序列化秒杀活动数据失败", zap.Error(err), zap.String("Data", data))
		return domain.SecondKillEvent{}, err
	}

	return event, nil
}

// SetSecondKillEvent 将秒杀活动设置到缓存中
func (c *lotteryDrawCache) SetSecondKillEvent(ctx context.Context, event domain.SecondKillEvent) error {
	key := generateSecondKillEventKey(event.ID)

	// 序列化为 JSON
	data, err := json.Marshal(event)
	if err != nil {
		c.logger.Error("序列化秒杀活动数据失败", zap.Error(err), zap.Any("Event", event))
		return err
	}

	// 设置缓存数据并设置过期时间
	if err := c.client.Set(ctx, key, data, secondKillEventTTL).Err(); err != nil {
		c.logger.Error("设置秒杀活动缓存失败", zap.Error(err), zap.Int("ID", event.ID))
		return err
	}

	return nil
}

// DeleteSecondKillEvent 从缓存中删除指定ID的秒杀活动
func (c *lotteryDrawCache) DeleteSecondKillEvent(ctx context.Context, id int) error {
	key := generateSecondKillEventKey(id)

	// 删除缓存键
	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.Error("删除秒杀活动缓存失败", zap.Error(err), zap.Int("ID", id))
		return err
	}

	return nil
}

// GetSecondKillEventWithLock 在缓存未命中时，使用分布式锁从数据库获取并设置缓存
func (c *lotteryDrawCache) GetSecondKillEventWithLock(ctx context.Context, id int, fetchFromDB func() (domain.SecondKillEvent, error)) (domain.SecondKillEvent, error) {
	// 尝试从缓存中获取数据
	event, err := c.GetSecondKillEvent(ctx, id)
	if err == nil {
		return event, nil
	}
	if !errors.Is(err, ErrCacheMiss) {
		return domain.SecondKillEvent{}, err
	}

	// 缓存未命中，尝试获取分布式锁
	lockKey := generateLockKey(generateSecondKillEventKey(id))
	lockValue := uuid.New().String()

	// 尝试获取锁
	acquired, err := c.acquireLock(ctx, lockKey, lockValue, lockTTL)
	if err != nil {
		c.logger.Error("获取分布式锁失败", zap.Error(err), zap.String("LockKey", lockKey))
		return domain.SecondKillEvent{}, err
	}

	if !acquired {
		// 获取锁失败，可能其他请求正在设置缓存，等待并重试获取缓存
		time.Sleep(lockRetryDelay)
		return c.GetSecondKillEvent(ctx, id)
	}

	// 确保在函数退出时释放锁
	defer func() {
		if err := c.releaseLock(ctx, lockKey, lockValue); err != nil {
			c.logger.Error("释放分布式锁失败", zap.Error(err), zap.String("LockKey", lockKey))
		}
	}()

	// 再次尝试从缓存中获取数据，以防止锁获取后，其他请求已经设置了缓存
	event, err = c.GetSecondKillEvent(ctx, id)
	if err == nil {
		return event, nil
	}
	if !errors.Is(err, ErrCacheMiss) {
		return domain.SecondKillEvent{}, err
	}

	// 从数据库获取数据
	event, err = fetchFromDB()
	if err != nil {
		return domain.SecondKillEvent{}, err
	}

	// 设置缓存
	if err := c.SetSecondKillEvent(ctx, event); err != nil {
		c.logger.Error("设置秒杀活动缓存失败", zap.Error(err), zap.Int("ID", id))
	}

	return event, nil
}

// acquireLock 尝试获取分布式锁
func (c *lotteryDrawCache) acquireLock(ctx context.Context, lockKey, lockValue string, ttl time.Duration) (bool, error) {
	// 使用 SET 命令获取锁
	ok, err := c.client.SetNX(ctx, lockKey, lockValue, ttl).Result()
	if err != nil {
		return false, err
	}

	return ok, nil
}

// releaseLock 释放分布式锁，使用 Lua 脚本确保操作的原子性
func (c *lotteryDrawCache) releaseLock(ctx context.Context, lockKey, lockValue string) error {
	// Lua 脚本：如果 key 的值等于 lockValue，则删除 key
	script := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
	`

	result, err := c.client.Eval(ctx, script, []string{lockKey}, lockValue).Result()
	if err != nil {
		return err
	}

	if result.(int64) == 0 {
		return errors.New("未能释放锁，可能锁已过期或被其他人释放")
	}

	return nil
}

// generateLotteryDrawKey 生成抽奖活动的缓存键
func generateLotteryDrawKey(id int) string {
	return fmt.Sprintf("%s%d", lotteryDrawKeyPrefix, id)
}

// generateSecondKillEventKey 生成秒杀活动的缓存键
func generateSecondKillEventKey(id int) string {
	return fmt.Sprintf("%s%d", secondKillEventKeyPrefix, id)
}

// generateLockKey 生成分布式锁的键
func generateLockKey(key string) string {
	return fmt.Sprintf("%s%s", lockKeyPrefix, key)
}
