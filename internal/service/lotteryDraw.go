package service

import (
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
)

// LotteryDrawService 定义了处理抽奖和秒杀相关业务逻辑的方法
type LotteryDrawService interface {
	// GetAllLotteryDraws 获取所有抽奖活动
	GetAllLotteryDraws() ([]domain.LotteryDraw, error)
	// CreateLotteryDraw 创建一个新的抽奖活动
	CreateLotteryDraw(input domain.CreateLotteryDrawInput) (domain.LotteryDraw, error)
	// GetLotteryDrawByID 根据ID获取指定的抽奖活动
	GetLotteryDrawByID(id string) (domain.LotteryDraw, error)
	// ParticipateLotteryDraw 参与指定ID的抽奖活动
	ParticipateLotteryDraw(id string, userID string) error

	// GetAllSecondKillEvents 获取所有秒杀活动
	GetAllSecondKillEvents() ([]domain.SecondKillEvent, error)
	// CreateSecondKillEvent 创建一个新的秒杀活动
	CreateSecondKillEvent(input domain.CreateSecondKillEventInput) (domain.SecondKillEvent, error)
	// GetSecondKillEventByID 根据ID获取指定的秒杀活动
	GetSecondKillEventByID(id string) (domain.SecondKillEvent, error)
	// ParticipateSecondKill 参与指定ID的秒杀活动
	ParticipateSecondKill(id string, userID string) error
}

type lotteryDrawService struct {
	repo repository.LotteryDrawRepository
}

func NewLotteryDrawService(repo repository.LotteryDrawRepository) LotteryDrawService {
	return &lotteryDrawService{
		repo: repo,
	}
}

func (l *lotteryDrawService) GetAllLotteryDraws() ([]domain.LotteryDraw, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawService) CreateLotteryDraw(input domain.CreateLotteryDrawInput) (domain.LotteryDraw, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawService) GetLotteryDrawByID(id string) (domain.LotteryDraw, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawService) ParticipateLotteryDraw(id string, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawService) GetAllSecondKillEvents() ([]domain.SecondKillEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawService) CreateSecondKillEvent(input domain.CreateSecondKillEventInput) (domain.SecondKillEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawService) GetSecondKillEventByID(id string) (domain.SecondKillEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawService) ParticipateSecondKill(id string, userID string) error {
	//TODO implement me
	panic("implement me")
}
