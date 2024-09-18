package repository

import (
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
)

type LotteryDrawRepository interface {
	// 抽奖相关的方法

	// GetAllLotteryDraws 获取所有抽奖活动
	GetAllLotteryDraws() ([]domain.LotteryDraw, error)

	// CreateLotteryDraw 创建一个新的抽奖活动
	CreateLotteryDraw(draw domain.LotteryDraw) (domain.LotteryDraw, error)

	// GetLotteryDrawByID 根据ID获取指定的抽奖活动
	GetLotteryDrawByID(id string) (domain.LotteryDraw, error)

	// ParticipateInLotteryDraw 参与抽奖活动
	ParticipateInLotteryDraw(participant domain.Participant) error

	// 秒杀相关的方法

	// GetAllSecondKillEvents 获取所有秒杀活动
	GetAllSecondKillEvents() ([]domain.SecondKillEvent, error)

	// CreateSecondKillEvent 创建一个新的秒杀活动
	CreateSecondKillEvent(event domain.SecondKillEvent) (domain.SecondKillEvent, error)

	// GetSecondKillEventByID 根据ID获取指定的秒杀活动
	GetSecondKillEventByID(id string) (domain.SecondKillEvent, error)

	// ParticipateInSecondKill 参与秒杀活动
	ParticipateInSecondKill(participant domain.Participant) error
}

type lotteryDrawRepository struct {
	dao dao.LotteryDrawDAO
}

func NewLotteryDrawRepository(dao dao.LotteryDrawDAO) LotteryDrawRepository {
	return &lotteryDrawRepository{
		dao: dao,
	}
}

func (l *lotteryDrawRepository) GetAllLotteryDraws() ([]domain.LotteryDraw, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawRepository) CreateLotteryDraw(draw domain.LotteryDraw) (domain.LotteryDraw, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawRepository) GetLotteryDrawByID(id string) (domain.LotteryDraw, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawRepository) ParticipateInLotteryDraw(participant domain.Participant) error {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawRepository) GetAllSecondKillEvents() ([]domain.SecondKillEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawRepository) CreateSecondKillEvent(event domain.SecondKillEvent) (domain.SecondKillEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawRepository) GetSecondKillEventByID(id string) (domain.SecondKillEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawRepository) ParticipateInSecondKill(participant domain.Participant) error {
	//TODO implement me
	panic("implement me")
}
