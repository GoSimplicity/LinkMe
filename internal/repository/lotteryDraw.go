package repository

import (
	"context"
	"errors"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

type LotteryDrawRepository interface {
	// 抽奖活动相关方法
	ListLotteryDraws(ctx context.Context, status string, pagination domain.Pagination) ([]domain.LotteryDraw, error)
	CreateLotteryDraw(ctx context.Context, draw domain.LotteryDraw) error
	GetLotteryDrawByID(ctx context.Context, id int) (domain.LotteryDraw, error)
	ExistsLotteryDrawByName(ctx context.Context, name string) (bool, error)
	HasUserParticipatedInLottery(ctx context.Context, id int, userID int64) (bool, error)
	AddLotteryParticipant(ctx context.Context, dp domain.Participant) error

	// 秒杀活动相关方法
	ListSecondKillEvents(ctx context.Context, status string, pagination domain.Pagination) ([]domain.SecondKillEvent, error)
	CreateSecondKillEvent(ctx context.Context, input domain.SecondKillEvent) error
	GetSecondKillEventByID(ctx context.Context, id int) (domain.SecondKillEvent, error)
	ExistsSecondKillEventByName(ctx context.Context, name string) (bool, error)
	HasUserParticipatedInSecondKill(ctx context.Context, id int, userID int64) (bool, error)
	AddSecondKillParticipant(ctx context.Context, dp domain.Participant) error

	// 活动状态管理方法
	ListPendingLotteryDraws(ctx context.Context, currentTime int64) ([]domain.LotteryDraw, error)
	UpdateLotteryDrawStatus(ctx context.Context, id int, status string) error
	ListPendingSecondKillEvents(ctx context.Context, currentTime int64) ([]domain.SecondKillEvent, error)
	UpdateSecondKillEventStatus(ctx context.Context, id int, status string) error
	ListActiveLotteryDraws(ctx context.Context, currentTime int64) ([]domain.LotteryDraw, error)
	ListActiveSecondKillEvents(ctx context.Context, currentTime int64) ([]domain.SecondKillEvent, error)
}

type lotteryDrawRepository struct {
	dao    dao.LotteryDrawDAO
	logger *zap.Logger
}

func NewLotteryDrawRepository(dao dao.LotteryDrawDAO, logger *zap.Logger) LotteryDrawRepository {
	return &lotteryDrawRepository{
		dao:    dao,
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

// CreateLotteryDraw 创建一个新的抽奖活动
func (r *lotteryDrawRepository) CreateLotteryDraw(ctx context.Context, draw domain.LotteryDraw) error {
	err := r.dao.CreateLotteryDraw(ctx, convertToDAOLotteryDraw(draw))
	if err != nil {
		r.logger.Error("创建抽奖活动失败", zap.Error(err), zap.String("name", draw.Name))
		return err
	}

	r.logger.Info("成功创建抽奖活动", zap.String("name", draw.Name))
	return nil
}

// GetLotteryDrawByID 根据 ID 获取指定的抽奖活动
func (r *lotteryDrawRepository) GetLotteryDrawByID(ctx context.Context, id int) (domain.LotteryDraw, error) {
	dbDraw, err := r.dao.GetLotteryDrawByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			r.logger.Warn("抽奖活动未找到", zap.Int("ID", id))
			return domain.LotteryDraw{}, err
		}
		r.logger.Error("获取抽奖活动失败", zap.Error(err), zap.Int("ID", id))
		return domain.LotteryDraw{}, err
	}

	return convertToDomainLotteryDraw(dbDraw), nil
}

// ExistsLotteryDrawByName 检查抽奖活动名称是否存在
func (r *lotteryDrawRepository) ExistsLotteryDrawByName(ctx context.Context, name string) (bool, error) {
	exists, err := r.dao.ExistsLotteryDrawByName(ctx, name)
	if err != nil {
		r.logger.Error("检查抽奖活动名称是否存在失败", zap.Error(err), zap.String("name", name))
		return false, err
	}
	return exists, nil
}

// HasUserParticipatedInLottery 检查用户是否已参与过该抽奖活动
func (r *lotteryDrawRepository) HasUserParticipatedInLottery(ctx context.Context, id int, userID int64) (bool, error) {
	participated, err := r.dao.HasUserParticipatedInLottery(ctx, id, userID)
	if err != nil {
		r.logger.Error("检查用户是否参与过抽奖活动失败", zap.Error(err), zap.Int("id", id), zap.Int64("userID", userID))
		return false, err
	}

	return participated, nil
}

// AddLotteryParticipant 添加用户抽奖参与记录
func (r *lotteryDrawRepository) AddLotteryParticipant(ctx context.Context, dp domain.Participant) error {
	err := r.dao.AddParticipant(ctx, convertToDAOParticipant(dp))
	if err != nil {
		r.logger.Error("添加抽奖参与者失败", zap.Error(err), zap.Int("LotteryID", *dp.LotteryID), zap.Int64("UserID", dp.UserID))
		return err
	}

	r.logger.Info("成功添加抽奖参与者", zap.Int("LotteryID", *dp.LotteryID), zap.Int64("UserID", dp.UserID))
	return nil
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

// CreateSecondKillEvent 创建一个新的秒杀活动
func (r *lotteryDrawRepository) CreateSecondKillEvent(ctx context.Context, input domain.SecondKillEvent) error {
	err := r.dao.CreateSecondKillEvent(ctx, convertToDAOSecondKillEvent(input))
	if err != nil {
		r.logger.Error("创建秒杀活动失败", zap.Error(err), zap.String("name", input.Name))
		return err
	}

	r.logger.Info("成功创建秒杀活动", zap.String("name", input.Name))
	return nil
}

// GetSecondKillEventByID 根据 ID 获取指定的秒杀活动
func (r *lotteryDrawRepository) GetSecondKillEventByID(ctx context.Context, id int) (domain.SecondKillEvent, error) {
	dbEvent, err := r.dao.GetSecondKillEventByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			r.logger.Warn("秒杀活动未找到", zap.Int("ID", id))
			return domain.SecondKillEvent{}, err
		}
		r.logger.Error("获取秒杀活动失败", zap.Error(err), zap.Int("ID", id))
		return domain.SecondKillEvent{}, err
	}

	return convertToDomainSecondKillEvent(dbEvent), nil
}

// ExistsSecondKillEventByName 检查秒杀活动名称是否存在
func (r *lotteryDrawRepository) ExistsSecondKillEventByName(ctx context.Context, name string) (bool, error) {
	exists, err := r.dao.ExistsSecondKillEventByName(ctx, name)
	if err != nil {
		r.logger.Error("检查秒杀活动名称是否存在失败", zap.Error(err), zap.String("name", name))
		return false, err
	}
	return exists, nil
}

// HasUserParticipatedInSecondKill 检查用户是否已参与过该秒杀活动
func (r *lotteryDrawRepository) HasUserParticipatedInSecondKill(ctx context.Context, id int, userID int64) (bool, error) {
	participated, err := r.dao.HasUserParticipatedInSecondKill(ctx, id, userID)
	if err != nil {
		r.logger.Error("检查用户是否参与过秒杀活动失败", zap.Error(err), zap.Int("id", id), zap.Int64("userID", userID))
		return false, err
	}
	return participated, nil
}

// AddSecondKillParticipant 添加用户秒杀参与记录
func (r *lotteryDrawRepository) AddSecondKillParticipant(ctx context.Context, dp domain.Participant) error {
	err := r.dao.AddParticipant(ctx, convertToDAOParticipant(dp))
	if err != nil {
		r.logger.Error("添加秒杀参与者失败", zap.Error(err), zap.Int("SecondKillID", *dp.SecondKillID), zap.Int64("UserID", dp.UserID))
		return err
	}

	r.logger.Info("成功添加秒杀参与者", zap.Int("SecondKillID", *dp.SecondKillID), zap.Int64("UserID", dp.UserID))
	return nil
}

// ListPendingLotteryDraws 获取所有待激活的抽奖活动
func (r *lotteryDrawRepository) ListPendingLotteryDraws(ctx context.Context, currentTime int64) ([]domain.LotteryDraw, error) {
	lotteryDraws, err := r.dao.ListPendingLotteryDraws(ctx, currentTime)
	if err != nil {
		r.logger.Error("获取待激活抽奖活动失败", zap.Error(err))
		return nil, err
	}

	return convertToDomainLotteryDraws(lotteryDraws), nil
}

// UpdateLotteryDrawStatus 更新抽奖活动的状态
func (r *lotteryDrawRepository) UpdateLotteryDrawStatus(ctx context.Context, id int, status string) error {
	err := r.dao.UpdateLotteryDrawStatus(ctx, id, status)
	if err != nil {
		r.logger.Error("更新抽奖活动状态失败", zap.Int("id", id), zap.String("status", status), zap.Error(err))
		return err
	}

	r.logger.Info("成功更新抽奖活动状态", zap.Int("id", id), zap.String("status", status))
	return nil
}

// ListPendingSecondKillEvents 获取所有待激活的秒杀活动
func (r *lotteryDrawRepository) ListPendingSecondKillEvents(ctx context.Context, currentTime int64) ([]domain.SecondKillEvent, error) {
	secondKillEvents, err := r.dao.ListPendingSecondKillEvents(ctx, currentTime)
	if err != nil {
		r.logger.Error("获取待激活秒杀活动失败", zap.Error(err))
		return nil, err
	}

	return convertToDomainSecondKillEvents(secondKillEvents), nil
}

// UpdateSecondKillEventStatus 更新秒杀活动的状态
func (r *lotteryDrawRepository) UpdateSecondKillEventStatus(ctx context.Context, id int, status string) error {
	err := r.dao.UpdateSecondKillEventStatus(ctx, id, status)
	if err != nil {
		r.logger.Error("更新秒杀活动状态失败", zap.Int("id", id), zap.String("status", status), zap.Error(err))
		return err
	}

	r.logger.Info("成功更新秒杀活动状态", zap.Int("id", id), zap.String("status", status))
	return nil
}

// ListActiveLotteryDraws 获取所有进行中的抽奖活动
func (r *lotteryDrawRepository) ListActiveLotteryDraws(ctx context.Context, currentTime int64) ([]domain.LotteryDraw, error) {
	lotteryDraws, err := r.dao.ListActiveLotteryDraws(ctx, currentTime)
	if err != nil {
		r.logger.Error("获取进行中的抽奖活动失败", zap.Error(err))
		return nil, err
	}

	return convertToDomainLotteryDraws(lotteryDraws), nil
}

// ListActiveSecondKillEvents 获取所有进行中的秒杀活动
func (r *lotteryDrawRepository) ListActiveSecondKillEvents(ctx context.Context, currentTime int64) ([]domain.SecondKillEvent, error) {
	secondKillEvents, err := r.dao.ListActiveSecondKillEvents(ctx, currentTime)
	if err != nil {
		r.logger.Error("获取进行中的秒杀活动失败", zap.Error(err))
		return nil, err
	}

	return convertToDomainSecondKillEvents(secondKillEvents), nil
}

// convertToDAOLotteryDraw 将 domain.LotteryDraw 转换为 dao.LotteryDraw
func convertToDAOLotteryDraw(d domain.LotteryDraw) dao.LotteryDraw {
	return dao.LotteryDraw{
		ID:           d.ID,
		Name:         d.Name,
		Description:  d.Description,
		StartTime:    d.StartTime,
		EndTime:      d.EndTime,
		Status:       d.Status,
		Participants: convertToDAOParticipants(d.Participants),
	}
}

// convertToDomainLotteryDraw 将 dao.LotteryDraw 转换为 domain.LotteryDraw
func convertToDomainLotteryDraw(d dao.LotteryDraw) domain.LotteryDraw {
	return domain.LotteryDraw{
		ID:           d.ID,
		Name:         d.Name,
		Description:  d.Description,
		StartTime:    d.StartTime,
		EndTime:      d.EndTime,
		Status:       d.Status,
		Participants: convertToDomainParticipants(d.Participants),
	}
}

// convertToDomainLotteryDraws 将 dao.LotteryDraw 列表转换为 domain.LotteryDraw 列表
func convertToDomainLotteryDraws(daoLotteryDraws []dao.LotteryDraw) []domain.LotteryDraw {
	domainLotteryDraws := make([]domain.LotteryDraw, 0, len(daoLotteryDraws))

	for _, d := range daoLotteryDraws {
		domainLotteryDraws = append(domainLotteryDraws, convertToDomainLotteryDraw(d))
	}

	return domainLotteryDraws
}

// convertToDAOSecondKillEvent 将 domain.SecondKillEvent 转换为 dao.SecondKillEvent
func convertToDAOSecondKillEvent(e domain.SecondKillEvent) dao.SecondKillEvent {
	return dao.SecondKillEvent{
		ID:           e.ID,
		Name:         e.Name,
		Description:  e.Description,
		StartTime:    e.StartTime,
		EndTime:      e.EndTime,
		Status:       e.Status,
		Participants: convertToDAOParticipants(e.Participants),
	}
}

// convertToDomainSecondKillEvent 将 dao.SecondKillEvent 转换为 domain.SecondKillEvent
func convertToDomainSecondKillEvent(e dao.SecondKillEvent) domain.SecondKillEvent {
	return domain.SecondKillEvent{
		ID:           e.ID,
		Name:         e.Name,
		Description:  e.Description,
		StartTime:    e.StartTime,
		EndTime:      e.EndTime,
		Status:       e.Status,
		Participants: convertToDomainParticipants(e.Participants),
	}
}

// convertToDomainSecondKillEvents 将 dao.SecondKillEvent 列表转换为 domain.SecondKillEvent 列表
func convertToDomainSecondKillEvents(daoSecondKillEvents []dao.SecondKillEvent) []domain.SecondKillEvent {
	domainSecondKillEvents := make([]domain.SecondKillEvent, 0, len(daoSecondKillEvents))

	for _, e := range daoSecondKillEvents {
		domainSecondKillEvents = append(domainSecondKillEvents, convertToDomainSecondKillEvent(e))
	}

	return domainSecondKillEvents
}

// convertToDAOParticipant 将 domain.Participant 转换为 dao.Participant
func convertToDAOParticipant(p domain.Participant) dao.Participant {
	return dao.Participant{
		ID:             p.ID,
		LotteryID:      p.LotteryID,
		SecondKillID:   p.SecondKillID,
		UserID:         p.UserID,
		ParticipatedAt: p.ParticipatedAt,
	}
}

// convertToDAOParticipants 将 domain.Participant 列表转换为 dao.Participant 列表
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

// convertToDomainParticipant 将 dao.Participant 转换为 domain.Participant
func convertToDomainParticipant(p dao.Participant) domain.Participant {
	return domain.Participant{
		ID:             p.ID,
		LotteryID:      p.LotteryID,
		SecondKillID:   p.SecondKillID,
		UserID:         p.UserID,
		ParticipatedAt: p.ParticipatedAt,
	}
}

// convertToDomainParticipants 将 dao.Participant 列表转换为 domain.Participant 列表
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
