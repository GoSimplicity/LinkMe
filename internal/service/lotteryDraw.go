package service

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/sync/semaphore"
)

type LotteryDrawService interface {
	// 抽奖活动相关方法
	ListLotteryDraws(ctx context.Context, status string, pagination domain.Pagination) ([]domain.LotteryDraw, error)
	CreateLotteryDraw(ctx context.Context, input domain.LotteryDraw) error
	GetLotteryDrawByID(ctx context.Context, id int) (domain.LotteryDraw, error)
	ParticipateLotteryDraw(ctx context.Context, id int, userID int64) error

	// 秒杀活动相关方法
	ListSecondKillEvents(ctx context.Context, status string, pagination domain.Pagination) ([]domain.SecondKillEvent, error)
	CreateSecondKillEvent(ctx context.Context, input domain.SecondKillEvent) error
	GetSecondKillEventByID(ctx context.Context, id int) (domain.SecondKillEvent, error)
	ParticipateSecondKill(ctx context.Context, id int, userID int64) error

	// 关闭服务
	Close() error
}

type lotteryDrawService struct {
	repo repository.LotteryDrawRepository
	l    *zap.Logger

	// 并发控制
	lotterySem    *semaphore.Weighted
	secondKillSem *semaphore.Weighted

	// 后台任务管理
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	closeOnce sync.Once

	// 活动锁管理
	locks sync.Map // map[int]*sync.RWMutex
}

const (
	// 定义抽奖和秒杀的最大并发数
	maxLotteryConcurrency    = 1000
	maxSecondKillConcurrency = 1000
)

func NewLotteryDrawService(repo repository.LotteryDrawRepository, l *zap.Logger) LotteryDrawService {
	ctx, cancel := context.WithCancel(context.Background())
	service := &lotteryDrawService{
		repo:          repo,
		l:             l,
		lotterySem:    semaphore.NewWeighted(maxLotteryConcurrency),
		secondKillSem: semaphore.NewWeighted(maxSecondKillConcurrency),
		ctx:           ctx,
		cancel:        cancel,
	}

	// 启动后台任务更新活动状态
	service.wg.Add(1)
	go service.runStatusUpdater()

	return service
}

// runStatusUpdater 启动后台定时任务，每10秒更新一次活动状态
func (s *lotteryDrawService) runStatusUpdater() {
	defer s.wg.Done()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			s.l.Info("状态更新器已停止")
			return
		case <-ticker.C:
			s.updateEventStatuses()
		}
	}
}

// updateEventStatuses 更新所有活动的状态
func (s *lotteryDrawService) updateEventStatuses() {
	currentTime := time.Now().Unix()

	// 更新抽奖活动状态
	if err := s.updateLotteryDrawStatuses(currentTime); err != nil {
		s.l.Error("更新抽奖活动状态失败", zap.Error(err))
	}

	// 更新秒杀活动状态
	if err := s.updateSecondKillEventStatuses(currentTime); err != nil {
		s.l.Error("更新秒杀活动状态失败", zap.Error(err))
	}
}

// updateLotteryDrawStatuses 更新所有抽奖活动的状态
func (s *lotteryDrawService) updateLotteryDrawStatuses(currentTime int64) error {
	ctx := context.Background()

	// 获取所有待开始的抽奖活动
	pendingDraws, err := s.repo.ListPendingLotteryDraws(ctx, currentTime)
	if err != nil {
		return err
	}

	for _, draw := range pendingDraws {
		// 如果活动开始时间已到，更新状态为进行中
		if draw.StartTime <= currentTime {
			if err := s.repo.UpdateLotteryDrawStatus(ctx, draw.ID, domain.LotteryStatusActive); err != nil {
				s.l.Error("更新抽奖活动状态为进行中失败",
					zap.Int("id", draw.ID),
					zap.Error(err))
			} else {
				s.l.Info("抽奖活动状态已更新为进行中",
					zap.Int("id", draw.ID))
			}
		}
	}

	// 获取所有进行中的抽奖活动
	activeDraws, err := s.repo.ListActiveLotteryDraws(ctx, currentTime)
	if err != nil {
		return err
	}

	for _, draw := range activeDraws {
		// 如果活动结束时间已到，更新状态为已完成
		if draw.EndTime <= currentTime {
			if err := s.repo.UpdateLotteryDrawStatus(ctx, draw.ID, domain.LotteryStatusCompleted); err != nil {
				s.l.Error("更新抽奖活动状态为已完成失败",
					zap.Int("id", draw.ID),
					zap.Error(err))
			} else {
				s.l.Info("抽奖活动状态已更新为已完成",
					zap.Int("id", draw.ID))
			}
		}
	}

	return nil
}

// updateSecondKillEventStatuses 更新所有秒杀活动的状态
func (s *lotteryDrawService) updateSecondKillEventStatuses(currentTime int64) error {
	ctx := context.Background()

	// 获取所有待开始的秒杀活动
	pendingEvents, err := s.repo.ListPendingSecondKillEvents(ctx, currentTime)
	if err != nil {
		return err
	}

	for _, event := range pendingEvents {
		// 如果活动开始时间已到，更新状态为进行中
		if event.StartTime <= currentTime {
			if err := s.repo.UpdateSecondKillEventStatus(ctx, event.ID, domain.SecondKillStatusActive); err != nil {
				s.l.Error("更新秒杀活动状态为进行中失败",
					zap.Int("id", event.ID),
					zap.Error(err))
			} else {
				s.l.Info("秒杀活动状态已更新为进行中",
					zap.Int("id", event.ID))
			}
		}
	}

	// 获取所有进行中的秒杀活动
	activeEvents, err := s.repo.ListActiveSecondKillEvents(ctx, currentTime)
	if err != nil {
		return err
	}

	for _, event := range activeEvents {
		// 如果活动结束时间已到，更新状态为已完成
		if event.EndTime <= currentTime {
			if err := s.repo.UpdateSecondKillEventStatus(ctx, event.ID, domain.SecondKillStatusCompleted); err != nil {
				s.l.Error("更新秒杀活动状态为已完成失败",
					zap.Int("id", event.ID),
					zap.Error(err))
			} else {
				s.l.Info("秒杀活动状态已更新为已完成",
					zap.Int("id", event.ID))
			}
		}
	}

	return nil
}

// Close 关闭服务，停止后台任务
func (s *lotteryDrawService) Close() error {
	s.closeOnce.Do(func() {
		// 取消上下文，通知后台任务停止
		s.cancel()

		// 等待后台任务完成
		s.wg.Wait()

		// 记录关闭信息
		s.l.Info("lottery draw service closed")
	})

	return nil
}

// ParticipateLotteryDraw 允许用户参与抽奖活动
func (s *lotteryDrawService) ParticipateLotteryDraw(ctx context.Context, id int, userID int64) error {
	// 获取信号量，限制并发
	if err := s.lotterySem.Acquire(ctx, 1); err != nil {
		s.l.Error("failed to acquire lottery semaphore", zap.Error(err))
		return err
	}
	defer s.lotterySem.Release(1)

	return s.processLotteryParticipation(ctx, id, userID)
}

// ParticipateSecondKill 允许用户参与秒杀活动
func (s *lotteryDrawService) ParticipateSecondKill(ctx context.Context, id int, userID int64) error {
	// 获取信号量，限制并发
	if err := s.secondKillSem.Acquire(ctx, 1); err != nil {
		s.l.Error("failed to acquire second kill semaphore", zap.Error(err))
		return err
	}
	defer s.secondKillSem.Release(1)

	return s.processSecondKillParticipation(ctx, id, userID)
}

// processLotteryParticipation 处理抽奖参与逻辑
func (s *lotteryDrawService) processLotteryParticipation(ctx context.Context, id int, userID int64) error {
	// 获取活动的读写锁
	lock := s.getLock(id)
	lock.Lock()
	defer lock.Unlock()

	currentTime := time.Now().Unix()

	// 更新活动状态（确保实时性）
	if err := s.updateSingleLotteryDrawStatus(ctx, id, currentTime); err != nil {
		s.l.Error("更新单个抽奖活动状态失败", zap.Int("id", id), zap.Error(err))
		return err
	}

	// 获取并验证抽奖活动
	lotteryDraw, err := s.repo.GetLotteryDrawByID(ctx, id)
	if err != nil {
		s.l.Error("failed to get lottery draw by ID", zap.Int("id", id), zap.Error(err))
		return err
	}

	// 验证活动状态
	if err := s.validateLotteryDraw(lotteryDraw, currentTime); err != nil {
		s.l.Error("lottery draw validation failed", zap.Int("id", id), zap.Error(err))
		return err
	}

	// 检查用户是否已参与
	alreadyParticipated, err := s.repo.HasUserParticipatedInLottery(ctx, id, userID)
	if err != nil {
		s.l.Error("failed to check user participation in lottery", zap.Int("id", id), zap.Int64("userID", userID), zap.Error(err))
		return err
	}

	if alreadyParticipated {
		return errors.New("用户已参与此抽奖活动")
	}

	// 创建参与记录
	participant := domain.Participant{
		ID:             generateUUID(),
		LotteryID:      &id,
		UserID:         userID,
		ParticipatedAt: currentTime,
	}

	// 添加参与者
	if err := s.repo.AddLotteryParticipant(ctx, participant); err != nil {
		s.l.Error("failed to add lottery participant", zap.Int("id", id), zap.Int64("userID", userID), zap.Error(err))
		return err
	}

	s.l.Info("user participated in lottery draw", zap.Int("id", id), zap.Int64("userID", userID))

	return nil
}

// processSecondKillParticipation 处理秒杀参与逻辑
func (s *lotteryDrawService) processSecondKillParticipation(ctx context.Context, id int, userID int64) error {
	// 获取活动的读写锁
	lock := s.getLock(id)
	lock.Lock()
	defer lock.Unlock()

	currentTime := time.Now().Unix()

	// 更新活动状态（确保实时性）
	if err := s.updateSingleSecondKillEventStatus(ctx, id, currentTime); err != nil {
		s.l.Error("更新单个秒杀活动状态失败", zap.Int("id", id), zap.Error(err))
		return err
	}

	// 获取并验证秒杀活动
	secondKillEvent, err := s.repo.GetSecondKillEventByID(ctx, id)
	if err != nil {
		s.l.Error("failed to get second kill event by ID", zap.Int("id", id), zap.Error(err))
		return err
	}

	// 验证活动状态
	if err := s.validateSecondKillEvent(secondKillEvent, currentTime); err != nil {
		s.l.Error("second kill event validation failed", zap.Int("id", id), zap.Error(err))
		return err
	}

	// 检查用户是否已参与
	alreadyParticipated, err := s.repo.HasUserParticipatedInSecondKill(ctx, id, userID)
	if err != nil {
		s.l.Error("failed to check user participation in second kill", zap.Int("id", id), zap.Int64("userID", userID), zap.Error(err))
		return err
	}

	if alreadyParticipated {
		return errors.New("用户已参与此秒杀活动")
	}

	// 创建参与记录
	participant := domain.Participant{
		ID:             generateUUID(),
		SecondKillID:   &id,
		UserID:         userID,
		ParticipatedAt: currentTime,
	}

	// 添加参与者
	if err := s.repo.AddSecondKillParticipant(ctx, participant); err != nil {
		s.l.Error("failed to add second kill participant", zap.Int("id", id), zap.Int64("userID", userID), zap.Error(err))
		return err
	}

	s.l.Info("user participated in second kill event", zap.Int("id", id), zap.Int64("userID", userID))

	return nil
}

// updateSingleLotteryDrawStatus 确保单个抽奖活动的状态是最新的
func (s *lotteryDrawService) updateSingleLotteryDrawStatus(ctx context.Context, id int, currentTime int64) error {
	// 获取抽奖活动
	lotteryDraw, err := s.repo.GetLotteryDrawByID(ctx, id)
	if err != nil {
		return err
	}

	// 更新状态到当前时间
	switch lotteryDraw.Status {
	case domain.LotteryStatusPending:
		if lotteryDraw.StartTime <= currentTime {
			if err := s.repo.UpdateLotteryDrawStatus(ctx, id, domain.LotteryStatusActive); err != nil {
				return err
			}
			s.l.Info("抽奖活动状态已更新为进行中", zap.Int("id", id))
		}
	case domain.LotteryStatusActive:
		if lotteryDraw.EndTime <= currentTime {
			if err := s.repo.UpdateLotteryDrawStatus(ctx, id, domain.LotteryStatusCompleted); err != nil {
				return err
			}
			s.l.Info("抽奖活动状态已更新为已完成", zap.Int("id", id))
		}
	}

	return nil
}

// updateSingleSecondKillEventStatus 确保单个秒杀活动的状态是最新的
func (s *lotteryDrawService) updateSingleSecondKillEventStatus(ctx context.Context, id int, currentTime int64) error {
	// 获取秒杀活动
	secondKillEvent, err := s.repo.GetSecondKillEventByID(ctx, id)
	if err != nil {
		return err
	}

	// 更新状态到当前时间
	switch secondKillEvent.Status {
	case domain.SecondKillStatusPending:
		if secondKillEvent.StartTime <= currentTime {
			if err := s.repo.UpdateSecondKillEventStatus(ctx, id, domain.SecondKillStatusActive); err != nil {
				return err
			}
			s.l.Info("秒杀活动状态已更新为进行中", zap.Int("id", id))
		}
	case domain.SecondKillStatusActive:
		if secondKillEvent.EndTime <= currentTime {
			if err := s.repo.UpdateSecondKillEventStatus(ctx, id, domain.SecondKillStatusCompleted); err != nil {
				return err
			}
			s.l.Info("秒杀活动状态已更新为已完成", zap.Int("id", id))
		}
	}

	return nil
}

// validateLotteryDraw 验证抽奖活动的状态和时间
func (s *lotteryDrawService) validateLotteryDraw(lotteryDraw domain.LotteryDraw, currentTime int64) error {
	if lotteryDraw.Status != domain.LotteryStatusActive {
		return errors.New("无法参与非进行中的抽奖活动")
	}

	if currentTime < lotteryDraw.StartTime || currentTime > lotteryDraw.EndTime {
		return errors.New("抽奖活动当前不在有效期内")
	}

	return nil
}

// validateSecondKillEvent 验证秒杀活动的状态和时间
func (s *lotteryDrawService) validateSecondKillEvent(event domain.SecondKillEvent, currentTime int64) error {
	if event.Status != domain.SecondKillStatusActive {
		return errors.New("无法参与非进行中的秒杀活动")
	}

	if currentTime < event.StartTime || currentTime > event.EndTime {
		return errors.New("秒杀活动当前不在有效期内")
	}

	return nil
}

// ListLotteryDraws 分页获取所有抽奖活动
func (s *lotteryDrawService) ListLotteryDraws(ctx context.Context, status string, pagination domain.Pagination) ([]domain.LotteryDraw, error) {
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset

	lotteries, err := s.repo.ListLotteryDraws(ctx, status, pagination)
	if err != nil {
		s.l.Error("failed to list lottery draws", zap.String("status", status), zap.Error(err))
		return nil, err
	}

	return lotteries, nil
}

// CreateLotteryDraw 创建新的抽奖活动
func (s *lotteryDrawService) CreateLotteryDraw(ctx context.Context, input domain.LotteryDraw) error {
	// 验证输入
	if err := validateLotteryDrawInput(input); err != nil {
		s.l.Error("invalid lottery draw input", zap.Error(err))
		return err
	}

	// 检查名称唯一性
	exists, err := s.repo.ExistsLotteryDrawByName(ctx, input.Name)
	if err != nil {
		s.l.Error("failed to check lottery draw name uniqueness", zap.String("name", input.Name), zap.Error(err))
		return err
	}

	if exists {
		return errors.New("同名的抽奖活动已存在")
	}

	// 设置状态
	currentTime := time.Now().Unix()
	var status string
	switch {
	case input.EndTime <= currentTime:
		status = domain.LotteryStatusCompleted
	case input.StartTime <= currentTime:
		status = domain.LotteryStatusActive
	default:
		status = domain.LotteryStatusPending
	}

	// 创建抽奖活动
	lotteryDraw := domain.LotteryDraw{
		Name:        input.Name,
		Description: input.Description,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Status:      status,
	}

	if err := s.repo.CreateLotteryDraw(ctx, lotteryDraw); err != nil {
		s.l.Error("failed to create lottery draw", zap.String("name", input.Name), zap.Error(err))
		return err
	}

	s.l.Info("lottery draw created", zap.String("name", input.Name))

	return nil
}

// GetLotteryDrawByID 根据ID获取抽奖活动
func (s *lotteryDrawService) GetLotteryDrawByID(ctx context.Context, id int) (domain.LotteryDraw, error) {
	lotteryDraw, err := s.repo.GetLotteryDrawByID(ctx, id)
	if err != nil {
		s.l.Error("failed to get lottery draw by ID", zap.Int("id", id), zap.Error(err))
		return domain.LotteryDraw{}, err
	}

	return lotteryDraw, nil
}

// ListSecondKillEvents 分页获取所有秒杀活动
func (s *lotteryDrawService) ListSecondKillEvents(ctx context.Context, status string, pagination domain.Pagination) ([]domain.SecondKillEvent, error) {
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset

	events, err := s.repo.ListSecondKillEvents(ctx, status, pagination)
	if err != nil {
		s.l.Error("failed to list second kill events", zap.String("status", status), zap.Error(err))
		return nil, err
	}

	return events, nil
}

// CreateSecondKillEvent 创建新的秒杀活动
func (s *lotteryDrawService) CreateSecondKillEvent(ctx context.Context, input domain.SecondKillEvent) error {
	// 验证输入
	if err := validateSecondKillEventInput(input); err != nil {
		s.l.Error("invalid second kill event input", zap.Error(err))
		return err
	}

	// 检查名称唯一性
	exists, err := s.repo.ExistsSecondKillEventByName(ctx, input.Name)
	if err != nil {
		s.l.Error("failed to check second kill event name uniqueness", zap.String("name", input.Name), zap.Error(err))
		return err
	}

	if exists {
		return errors.New("同名的秒杀活动已存在")
	}

	// 设置状态
	currentTime := time.Now().Unix()
	var status string

	switch {
	case input.EndTime <= currentTime:
		status = domain.SecondKillStatusCompleted
	case input.StartTime <= currentTime:
		status = domain.SecondKillStatusActive
	default:
		status = domain.SecondKillStatusPending
	}

	// 创建秒杀活动
	secondKillEvent := domain.SecondKillEvent{
		Name:        input.Name,
		Description: input.Description,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Status:      status,
	}

	if err := s.repo.CreateSecondKillEvent(ctx, secondKillEvent); err != nil {
		s.l.Error("failed to create second kill event", zap.String("name", input.Name), zap.Error(err))
		return err
	}

	s.l.Info("second kill event created", zap.String("name", input.Name))

	return nil
}

// GetSecondKillEventByID 根据ID获取秒杀活动
func (s *lotteryDrawService) GetSecondKillEventByID(ctx context.Context, id int) (domain.SecondKillEvent, error) {
	event, err := s.repo.GetSecondKillEventByID(ctx, id)
	if err != nil {
		s.l.Error("failed to get second kill event by ID", zap.Int("id", id), zap.Error(err))
		return domain.SecondKillEvent{}, err
	}

	return event, nil
}

// getLock 获取指定活动的读写锁，如果不存在则创建一个新的锁
func (s *lotteryDrawService) getLock(id int) *sync.RWMutex {
	actual, _ := s.locks.LoadOrStore(id, &sync.RWMutex{})
	return actual.(*sync.RWMutex)
}

// generateUUID 生成新的UUID
func generateUUID() string {
	return uuid.New().String()
}

// validateLotteryDrawInput 验证创建抽奖活动的输入
func validateLotteryDrawInput(input domain.LotteryDraw) error {
	if input.Name == "" {
		return errors.New("抽奖活动名称不能为空")
	}
	if input.StartTime >= input.EndTime {
		return errors.New("无效的抽奖活动时间范围")
	}
	return nil
}

// validateSecondKillEventInput 验证创建秒杀活动的输入
func validateSecondKillEventInput(input domain.SecondKillEvent) error {
	if input.Name == "" {
		return errors.New("秒杀活动名称不能为空")
	}
	if input.StartTime >= input.EndTime {
		return errors.New("无效的秒杀活动时间范围")
	}
	return nil
}
