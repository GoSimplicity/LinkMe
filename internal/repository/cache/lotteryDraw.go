package cache

import (
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/redis/go-redis/v9"
)

// LotteryDrawCache 定义了抽奖和秒杀活动的缓存接口
type LotteryDrawCache interface {
	// 抽奖相关的方法

	// GetLotteryDraw 从缓存中获取指定ID的抽奖活动
	GetLotteryDraw(id string) (*domain.LotteryDraw, error)

	// SetLotteryDraw 将抽奖活动设置到缓存中
	SetLotteryDraw(draw *domain.LotteryDraw) error

	// DeleteLotteryDraw 从缓存中删除指定ID的抽奖活动
	DeleteLotteryDraw(id string) error

	// 秒杀相关的方法

	// GetSecondKillEvent 从缓存中获取指定ID的秒杀活动
	GetSecondKillEvent(id string) (*domain.SecondKillEvent, error)

	// SetSecondKillEvent 将秒杀活动设置到缓存中
	SetSecondKillEvent(event *domain.SecondKillEvent) error

	// DeleteSecondKillEvent 从缓存中删除指定ID的秒杀活动
	DeleteSecondKillEvent(id string) error
}

type lotteryDrawCache struct {
	client redis.Cmdable
}

func NewLotteryDrawCache(client redis.Cmdable) LotteryDrawCache {
	return &lotteryDrawCache{
		client: client,
	}
}

func (l *lotteryDrawCache) GetLotteryDraw(id string) (*domain.LotteryDraw, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawCache) SetLotteryDraw(draw *domain.LotteryDraw) error {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawCache) DeleteLotteryDraw(id string) error {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawCache) GetSecondKillEvent(id string) (*domain.SecondKillEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawCache) SetSecondKillEvent(event *domain.SecondKillEvent) error {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawCache) DeleteSecondKillEvent(id string) error {
	//TODO implement me
	panic("implement me")
}
