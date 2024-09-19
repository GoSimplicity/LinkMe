package repository

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
)

type LotteryDrawRepository interface {
	// ListLotteryDraws 获取所有抽奖活动
	ListLotteryDraws(ctx context.Context, status string, pagination domain.Pagination) ([]domain.LotteryDraw, error)
	// CreateLotteryDraw 创建一个新的抽奖活动
	CreateLotteryDraw(ctx context.Context, draw domain.LotteryDraw) error
	// GetLotteryDrawByID 根据ID获取指定的抽奖活动
	GetLotteryDrawByID(ctx context.Context, id int) (domain.LotteryDraw, error)
	// ExistsLotteryDrawByName 检查抽奖活动名称是否存在
	ExistsLotteryDrawByName(ctx context.Context, name string) (bool, error)
	// HasUserParticipatedInLottery 检查用户是否已参与过该抽奖活动
	HasUserParticipatedInLottery(ctx context.Context, id int, userID int64) (bool, error)
	// AddLotteryParticipant 参与抽奖活动
	AddLotteryParticipant(ctx context.Context, dp domain.Participant) error

	// ListSecondKillEvents 获取所有秒杀活动
	ListSecondKillEvents(ctx context.Context, status string, pagination domain.Pagination) ([]domain.SecondKillEvent, error)
	// CreateSecondKillEvent 创建一个新的秒杀活动
	CreateSecondKillEvent(ctx context.Context, input domain.SecondKillEvent) error
	// GetSecondKillEventByID 根据ID获取指定的秒杀活动
	GetSecondKillEventByID(ctx context.Context, id int) (domain.SecondKillEvent, error)
	// ExistsSecondKillEventByName 检查秒杀活动名称是否存在
	ExistsSecondKillEventByName(ctx context.Context, name string) (bool, error)
	// HasUserParticipatedInSecondKill 检查用户是否已参与过该秒杀活动
	HasUserParticipatedInSecondKill(ctx context.Context, id int, userID int64) (bool, error)
	// AddSecondKillParticipant 参与秒杀活动
	AddSecondKillParticipant(ctx context.Context, dp domain.Participant) error
}

type lotteryDrawRepository struct {
	dao dao.LotteryDrawDAO
}

func NewLotteryDrawRepository(dao dao.LotteryDrawDAO) LotteryDrawRepository {
	return &lotteryDrawRepository{
		dao: dao,
	}
}

func (l lotteryDrawRepository) ListLotteryDraws(ctx context.Context, status string, pagination domain.Pagination) ([]domain.LotteryDraw, error) {
	//TODO implement me
	panic("implement me")
}

func (l lotteryDrawRepository) CreateLotteryDraw(ctx context.Context, draw domain.LotteryDraw) error {
	//TODO implement me
	panic("implement me")
}

func (l lotteryDrawRepository) GetLotteryDrawByID(ctx context.Context, id int) (domain.LotteryDraw, error) {
	//TODO implement me
	panic("implement me")
}

func (l lotteryDrawRepository) ExistsLotteryDrawByName(ctx context.Context, name string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (l lotteryDrawRepository) AddLotteryParticipant(ctx context.Context, dp domain.Participant) error {
	//TODO implement me
	panic("implement me")
}

func (l lotteryDrawRepository) ListSecondKillEvents(ctx context.Context, status string, pagination domain.Pagination) ([]domain.SecondKillEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (l lotteryDrawRepository) CreateSecondKillEvent(ctx context.Context, input domain.SecondKillEvent) error {
	//TODO implement me
	panic("implement me")
}

func (l lotteryDrawRepository) GetSecondKillEventByID(ctx context.Context, id int) (domain.SecondKillEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (l lotteryDrawRepository) HasUserParticipatedInLottery(ctx context.Context, id int, userID int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (l lotteryDrawRepository) ExistsSecondKillEventByName(ctx context.Context, name string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (l lotteryDrawRepository) HasUserParticipatedInSecondKill(ctx context.Context, id int, userID int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (l lotteryDrawRepository) AddSecondKillParticipant(ctx context.Context, dp domain.Participant) error {
	//TODO implement me
	panic("implement me")
}
