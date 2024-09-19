package service

import (
	"context"
	"errors"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/google/uuid"
)

type LotteryDrawService interface {
	// ListLotteryDraws 获取所有抽奖活动
	ListLotteryDraws(ctx context.Context, status string, pagination domain.Pagination) ([]domain.LotteryDraw, error)
	// CreateLotteryDraw 创建一个新的抽奖活动
	CreateLotteryDraw(ctx context.Context, input domain.LotteryDraw) error
	// GetLotteryDrawByID 根据ID获取指定的抽奖活动
	GetLotteryDrawByID(ctx context.Context, id int) (domain.LotteryDraw, error)
	// ParticipateLotteryDraw 参与指定ID的抽奖活动
	ParticipateLotteryDraw(ctx context.Context, id int, userID int64) error

	// ListSecondKillEvents 获取所有秒杀活动
	ListSecondKillEvents(ctx context.Context, status string, pagination domain.Pagination) ([]domain.SecondKillEvent, error)
	// CreateSecondKillEvent 创建一个新的秒杀活动
	CreateSecondKillEvent(ctx context.Context, input domain.SecondKillEvent) error
	// GetSecondKillEventByID 根据ID获取指定的秒杀活动
	GetSecondKillEventByID(ctx context.Context, id int) (domain.SecondKillEvent, error)
	// ParticipateSecondKill 参与指定ID的秒杀活动
	ParticipateSecondKill(ctx context.Context, id int, userID int64) error
}

type lotteryDrawService struct {
	repo repository.LotteryDrawRepository
}

func NewLotteryDrawService(repo repository.LotteryDrawRepository) LotteryDrawService {
	return &lotteryDrawService{
		repo: repo,
	}
}

// ListLotteryDraws 获取所有抽奖活动
func (l *lotteryDrawService) ListLotteryDraws(ctx context.Context, status string, pagination domain.Pagination) ([]domain.LotteryDraw, error) {
	lotteryDraws, err := l.repo.ListLotteryDraws(ctx, status, pagination)
	if err != nil {
		return nil, err
	}

	return lotteryDraws, nil
}

// CreateLotteryDraw 创建一个新的抽奖活动
func (l *lotteryDrawService) CreateLotteryDraw(ctx context.Context, input domain.LotteryDraw) error {
	// 检查活动名称是否唯一
	exists, err := l.repo.ExistsLotteryDrawByName(ctx, input.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("lottery draw with the same name already exists")
	}

	// 设置活动状态
	currentTime := time.Now().Unix()
	status := domain.LotteryStatusPending
	if input.StartTime <= currentTime {
		status = domain.LotteryStatusActive
	}

	// 创建新的抽奖活动对象
	lotteryDraw := domain.LotteryDraw{
		Name:        input.Name,
		Description: input.Description,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Status:      status,
	}

	// 创建抽奖活动
	err = l.repo.CreateLotteryDraw(ctx, lotteryDraw)
	if err != nil {
		return err
	}

	return nil
}

// GetLotteryDrawByID 根据ID获取指定的抽奖活动
func (l *lotteryDrawService) GetLotteryDrawByID(ctx context.Context, id int) (domain.LotteryDraw, error) {
	lotteryDraw, err := l.repo.GetLotteryDrawByID(ctx, id)
	if err != nil {
		return domain.LotteryDraw{}, err
	}

	return lotteryDraw, nil
}

// ParticipateLotteryDraw 参与指定ID的抽奖活动
func (l *lotteryDrawService) ParticipateLotteryDraw(ctx context.Context, id int, userID int64) error {
	// 获取当前时间
	currentTime := time.Now().Unix()

	// 获取抽奖活动
	lotteryDraw, err := l.repo.GetLotteryDrawByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查活动状态是否为进行中
	if lotteryDraw.Status != domain.LotteryStatusActive {
		return errors.New("cannot participate in a lottery draw that is not active")
	}

	// 检查活动是否在有效时间范围内
	if currentTime < lotteryDraw.StartTime || currentTime > lotteryDraw.EndTime {
		return errors.New("the lottery draw is not currently active")
	}

	// 检查用户是否已参与
	alreadyParticipated, err := l.repo.HasUserParticipatedInLottery(ctx, id, userID)
	if err != nil {
		return err
	}

	if alreadyParticipated {
		return errors.New("user has already participated in this lottery draw")
	}

	// 创建参与记录
	participant := domain.Participant{
		ID:             generateUUID(),
		ActivityID:     id,
		UserID:         userID,
		ParticipatedAt: currentTime,
	}

	// 调用仓库方法添加参与者
	err = l.repo.AddLotteryParticipant(ctx, participant)
	if err != nil {
		return err
	}

	return nil
}

// ListSecondKillEvents 获取所有秒杀活动
func (l *lotteryDrawService) ListSecondKillEvents(ctx context.Context, status string, pagination domain.Pagination) ([]domain.SecondKillEvent, error) {
	secondKillEvents, err := l.repo.ListSecondKillEvents(ctx, status, pagination)
	if err != nil {
		return nil, err
	}

	return secondKillEvents, nil
}

// CreateSecondKillEvent 创建一个新的秒杀活动
func (l *lotteryDrawService) CreateSecondKillEvent(ctx context.Context, input domain.SecondKillEvent) error {
	// 检查活动名称是否唯一
	exists, err := l.repo.ExistsSecondKillEventByName(ctx, input.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("second kill event with the same name already exists")
	}

	// 设置活动状态
	currentTime := time.Now().Unix()
	status := domain.SecondKillStatusPending
	if input.StartTime <= currentTime {
		status = domain.SecondKillStatusActive
	}

	// 创建新的秒杀活动对象
	secondKillEvent := domain.SecondKillEvent{
		Name:        input.Name,
		Description: input.Description,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Status:      status,
	}

	// 创建秒杀活动
	err = l.repo.CreateSecondKillEvent(ctx, secondKillEvent)
	if err != nil {
		return err
	}

	return nil
}

// GetSecondKillEventByID 根据ID获取指定的秒杀活动
func (l *lotteryDrawService) GetSecondKillEventByID(ctx context.Context, id int) (domain.SecondKillEvent, error) {
	secondKillEvent, err := l.repo.GetSecondKillEventByID(ctx, id)
	if err != nil {
		return domain.SecondKillEvent{}, err
	}

	return secondKillEvent, nil
}

// ParticipateSecondKill 参与指定ID的秒杀活动
func (l *lotteryDrawService) ParticipateSecondKill(ctx context.Context, id int, userID int64) error {
	// 获取当前时间
	currentTime := time.Now().Unix()

	// 获取秒杀活动
	secondKillEvent, err := l.repo.GetSecondKillEventByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查活动状态是否为进行中
	if secondKillEvent.Status != domain.SecondKillStatusActive {
		return errors.New("cannot participate in a second kill event that is not active")
	}

	// 检查活动是否在有效时间范围内
	if currentTime < secondKillEvent.StartTime || currentTime > secondKillEvent.EndTime {
		return errors.New("the second kill event is not currently active")
	}

	// 检查用户是否已参与
	alreadyParticipated, err := l.repo.HasUserParticipatedInSecondKill(ctx, id, userID)
	if err != nil {
		return err
	}
	if alreadyParticipated {
		return errors.New("user has already participated in this second kill event")
	}

	// 创建参与记录
	participant := domain.Participant{
		ID:             generateUUID(),
		ActivityID:     id,
		UserID:         userID,
		ParticipatedAt: currentTime,
	}

	// 调用仓库方法添加参与者
	err = l.repo.AddSecondKillParticipant(ctx, participant)
	if err != nil {
		return err
	}

	return nil
}

// generateUUID 生成一个新的UUID
func generateUUID() string {
	return uuid.New().String()
}
