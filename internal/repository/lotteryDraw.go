package repository

import (
	"context"
	"errors"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

type LotteryDrawRepository interface {
	ListLotteryDraws(ctx context.Context, status string, pagination domain.Pagination) ([]domain.LotteryDraw, error)
	CreateLotteryDraw(ctx context.Context, draw domain.LotteryDraw) error
	GetLotteryDrawByID(ctx context.Context, id int) (domain.LotteryDraw, error)
	ExistsLotteryDrawByName(ctx context.Context, name string) (bool, error)
	HasUserParticipatedInLottery(ctx context.Context, id int, userID int64) (bool, error)
	AddLotteryParticipant(ctx context.Context, dp domain.Participant) error

	ListSecondKillEvents(ctx context.Context, status string, pagination domain.Pagination) ([]domain.SecondKillEvent, error)
	CreateSecondKillEvent(ctx context.Context, input domain.SecondKillEvent) error
	GetSecondKillEventByID(ctx context.Context, id int) (domain.SecondKillEvent, error)
	ExistsSecondKillEventByName(ctx context.Context, name string) (bool, error)
	HasUserParticipatedInSecondKill(ctx context.Context, id int, userID int64) (bool, error)
	AddSecondKillParticipant(ctx context.Context, dp domain.Participant) error
}

type lotteryDrawRepository struct {
	dao    dao.LotteryDrawDAO
	cache  cache.LotteryDrawCache
	logger *zap.Logger
}

func NewLotteryDrawRepository(dao dao.LotteryDrawDAO, cache cache.LotteryDrawCache, logger *zap.Logger) LotteryDrawRepository {
	return &lotteryDrawRepository{
		dao:    dao,
		cache:  cache,
		logger: logger,
	}
}

// ListLotteryDraws 获取所有抽奖活动，支持状态过滤和分页
func (r *lotteryDrawRepository) ListLotteryDraws(ctx context.Context, status string, pagination domain.Pagination) ([]domain.LotteryDraw, error) {
	lotteryDraws, err := r.dao.ListLotteryDraws(ctx, status, pagination)
	if err != nil {
		r.logger.Error("获取抽奖活动列表失败", zap.Error(err))
		return nil, err
	}

	return convertToDomainLotteryDraws(lotteryDraws), nil
}

// CreateLotteryDraw 创建一个新的抽奖活动，并将其缓存
func (r *lotteryDrawRepository) CreateLotteryDraw(ctx context.Context, draw domain.LotteryDraw) error {
	err := r.dao.CreateLotteryDraw(ctx, convertToDAOLotteryDraw(draw))
	if err != nil {
		return err
	}

	// 更新缓存
	err = r.cache.SetLotteryDraw(ctx, draw)
	if err != nil {
		r.logger.Error("缓存创建的抽奖活动失败", zap.Error(err), zap.Int("ID", draw.ID))
	}

	return nil
}

// GetLotteryDrawByID 根据 ID 获取指定的抽奖活动，优先从缓存获取
func (r *lotteryDrawRepository) GetLotteryDrawByID(ctx context.Context, id int) (domain.LotteryDraw, error) {
	lotteryDraw, err := r.cache.GetLotteryDrawWithLock(ctx, id, func() (domain.LotteryDraw, error) {
		// 从数据库获取数据
		dbDraw, err := r.dao.GetLotteryDrawByID(ctx, id)
		if err != nil {
			return domain.LotteryDraw{}, err
		}
		domainDraw := convertToDomainLotteryDraw(dbDraw)
		return domainDraw, nil
	})

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			r.logger.Warn("抽奖活动未找到", zap.Int("ID", id))
			return domain.LotteryDraw{}, err
		}
		r.logger.Error("获取抽奖活动失败", zap.Error(err), zap.Int("ID", id))
		return domain.LotteryDraw{}, err
	}

	return lotteryDraw, nil
}

// ExistsLotteryDrawByName 检查抽奖活动名称是否存在
func (r *lotteryDrawRepository) ExistsLotteryDrawByName(ctx context.Context, name string) (bool, error) {
	return r.dao.ExistsLotteryDrawByName(ctx, name)
}

// HasUserParticipatedInLottery 检查用户是否已参与过该抽奖活动
func (r *lotteryDrawRepository) HasUserParticipatedInLottery(ctx context.Context, id int, userID int64) (bool, error) {
	return r.dao.HasUserParticipatedInLottery(ctx, id, userID)
}

// AddLotteryParticipant 添加用户抽奖参与记录，并更新缓存
func (r *lotteryDrawRepository) AddLotteryParticipant(ctx context.Context /**/, dp domain.Participant) error {
	err := r.dao.AddParticipant(ctx, convertToDAOParticipant(dp))
	if err != nil {
		return err
	}

	// 更新缓存：获取当前活动并更新缓存
	lotteryDraw, err := r.GetLotteryDrawByID(ctx, *dp.LotteryID)
	if err != nil {
		r.logger.Error("添加参与者后获取抽奖活动失败，无法更新缓存", zap.Error(err), zap.Int("ActivityID", *dp.LotteryID))
		return err
	}

	return r.cache.SetLotteryDraw(ctx, lotteryDraw)
}

// ListSecondKillEvents 获取所有秒杀活动，支持状态过滤和分页
func (r *lotteryDrawRepository) ListSecondKillEvents(ctx context.Context, status string, pagination domain.Pagination) ([]domain.SecondKillEvent, error) {
	secondKillEvents, err := r.dao.ListSecondKillEvents(ctx, status, pagination)
	if err != nil {
		r.logger.Error("获取秒杀活动列表失败", zap.Error(err))
		return nil, err
	}

	return convertToDomainSecondKillEvents(secondKillEvents), nil
}

// CreateSecondKillEvent 创建一个新的秒杀活动，并将其缓存
func (r *lotteryDrawRepository) CreateSecondKillEvent(ctx context.Context, input domain.SecondKillEvent) error {
	err := r.dao.CreateSecondKillEvent(ctx, convertToDAOSecondKillEvent(input))
	if err != nil {
		return err
	}

	// 更新缓存
	err = r.cache.SetSecondKillEvent(ctx, input)
	if err != nil {
		r.logger.Error("缓存创建的秒杀活动失败", zap.Error(err), zap.Int("ID", input.ID))
	}

	return nil
}

// GetSecondKillEventByID 根据 ID 获取指定的秒杀活动，优先从缓存获取
func (r *lotteryDrawRepository) GetSecondKillEventByID(ctx context.Context, id int) (domain.SecondKillEvent, error) {
	secondKillEvent, err := r.cache.GetSecondKillEventWithLock(ctx, id, func() (domain.SecondKillEvent, error) {
		// 从数据库获取数据
		dbEvent, err := r.dao.GetSecondKillEventByID(ctx, id)
		if err != nil {
			return domain.SecondKillEvent{}, err
		}
		domainEvent := convertToDomainSecondKillEvent(dbEvent)
		return domainEvent, nil
	})

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			r.logger.Warn("秒杀活动未找到", zap.Int("ID", id))
			return domain.SecondKillEvent{}, err
		}
		r.logger.Error("获取秒杀活动失败", zap.Error(err), zap.Int("ID", id))
		return domain.SecondKillEvent{}, err
	}

	return secondKillEvent, nil
}

// ExistsSecondKillEventByName 检查秒杀活动名称是否存在
func (r *lotteryDrawRepository) ExistsSecondKillEventByName(ctx context.Context, name string) (bool, error) {
	return r.dao.ExistsSecondKillEventByName(ctx, name)
}

// HasUserParticipatedInSecondKill 检查用户是否已参与过该秒杀活动
func (r *lotteryDrawRepository) HasUserParticipatedInSecondKill(ctx context.Context, id int, userID int64) (bool, error) {
	return r.dao.HasUserParticipatedInSecondKill(ctx, id, userID)
}

// AddSecondKillParticipant 添加用户秒杀参与记录，并更新缓存
func (r *lotteryDrawRepository) AddSecondKillParticipant(ctx context.Context, dp domain.Participant) error {
	err := r.dao.AddParticipant(ctx, convertToDAOParticipant(dp))
	if err != nil {
		return err
	}

	// 更新缓存：获取当前活动并更新缓存
	secondKillEvent, err := r.GetSecondKillEventByID(ctx, *dp.SecondKillID)
	if err != nil {
		r.logger.Error("添加参与者后获取秒杀活动失败，无法更新缓存", zap.Error(err), zap.Int("ActivityID", *dp.SecondKillID))
		return err
	}

	return r.cache.SetSecondKillEvent(ctx, secondKillEvent)
}

// ======================== 转换函数 ========================

// 将 domain.LotteryDraw 转换为 dao.LotteryDraw
func convertToDAOLotteryDraw(d domain.LotteryDraw) dao.LotteryDraw {
	return dao.LotteryDraw{
		ID:           d.ID,
		Name:         d.Name,
		Description:  d.Description,
		StartTime:    d.StartTime,
		EndTime:      d.EndTime,
		Status:       string(d.Status), // 将枚举状态转换为字符串
		Participants: convertToDAOParticipants(d.Participants),
	}
}

// 将 dao.LotteryDraw 转换为 domain.LotteryDraw
func convertToDomainLotteryDraw(d dao.LotteryDraw) domain.LotteryDraw {
	return domain.LotteryDraw{
		ID:           d.ID,
		Name:         d.Name,
		Description:  d.Description,
		StartTime:    d.StartTime,
		EndTime:      d.EndTime,
		Status:       domain.LotteryStatus(d.Status),
		Participants: convertToDomainParticipants(d.Participants),
	}
}

// 将 dao.LotteryDraw 列表转换为 domain.LotteryDraw 列表
func convertToDomainLotteryDraws(daoLotteryDraws []dao.LotteryDraw) []domain.LotteryDraw {
	domainLotteryDraws := make([]domain.LotteryDraw, 0, len(daoLotteryDraws))

	for _, d := range daoLotteryDraws {
		domainLotteryDraws = append(domainLotteryDraws, convertToDomainLotteryDraw(d))
	}

	return domainLotteryDraws
}

// 将 domain.SecondKillEvent 转换为 dao.SecondKillEvent
func convertToDAOSecondKillEvent(e domain.SecondKillEvent) dao.SecondKillEvent {
	return dao.SecondKillEvent{
		ID:           e.ID,
		Name:         e.Name,
		Description:  e.Description,
		StartTime:    e.StartTime,
		EndTime:      e.EndTime,
		Status:       string(e.Status),
		Participants: convertToDAOParticipants(e.Participants),
	}
}

// 将 dao.SecondKillEvent 转换为 domain.SecondKillEvent
func convertToDomainSecondKillEvent(e dao.SecondKillEvent) domain.SecondKillEvent {
	return domain.SecondKillEvent{
		ID:           e.ID,
		Name:         e.Name,
		Description:  e.Description,
		StartTime:    e.StartTime,
		EndTime:      e.EndTime,
		Status:       domain.SecondKillStatus(e.Status),
		Participants: convertToDomainParticipants(e.Participants),
	}
}

// 将 dao.SecondKillEvent 列表转换为 domain.SecondKillEvent 列表
func convertToDomainSecondKillEvents(daoSecondKillEvents []dao.SecondKillEvent) []domain.SecondKillEvent {
	domainSecondKillEvents := make([]domain.SecondKillEvent, 0, len(daoSecondKillEvents))

	for _, e := range daoSecondKillEvents {
		domainSecondKillEvents = append(domainSecondKillEvents, convertToDomainSecondKillEvent(e))
	}

	return domainSecondKillEvents
}

// 将 domain.Participant 转换为 dao.Participant
func convertToDAOParticipant(p domain.Participant) dao.Participant {
	return dao.Participant{
		ID:             p.ID,
		LotteryID:      p.LotteryID,
		SecondKillID:   p.SecondKillID,
		UserID:         p.UserID,
		ParticipatedAt: p.ParticipatedAt,
	}
}

// 将 domain.Participant 列表转换为 dao.Participant 列表
func convertToDAOParticipants(domainParticipants []domain.Participant) []dao.Participant {
	if len(domainParticipants) == 0 {
		return nil
	}

	daoParticipants := make([]dao.Participant, 0, len(domainParticipants))
	for _, p := range domainParticipants {
		daoParticipants = append(daoParticipants, convertToDAOParticipant(p))
	}

	return daoParticipants
}

// 将 dao.Participant 转换为 domain.Participant
func convertToDomainParticipant(p dao.Participant) domain.Participant {
	return domain.Participant{
		ID:             p.ID,
		LotteryID:      p.LotteryID,
		SecondKillID:   p.SecondKillID,
		UserID:         p.UserID,
		ParticipatedAt: p.ParticipatedAt,
	}
}

// 将 dao.Participant 列表转换为 domain.Participant 列表
func convertToDomainParticipants(daoParticipants []dao.Participant) []domain.Participant {
	if len(daoParticipants) == 0 {
		return nil
	}

	domainParticipants := make([]domain.Participant, 0, len(daoParticipants))

	for _, p := range daoParticipants {
		domainParticipants = append(domainParticipants, convertToDomainParticipant(p))
	}

	return domainParticipants
}
