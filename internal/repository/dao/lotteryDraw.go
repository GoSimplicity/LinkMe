package dao

import (
	"context"
	"errors"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LotteryDrawDAO interface {
	CreateLotteryDraw(ctx context.Context, model LotteryDraw) error
	GetLotteryDrawByID(ctx context.Context, id int) (LotteryDraw, error)
	ListLotteryDraws(ctx context.Context, status string, pagination domain.Pagination) ([]LotteryDraw, error)
	ExistsLotteryDrawByName(ctx context.Context, name string) (bool, error)
	HasUserParticipatedInLottery(ctx context.Context, id int, userID int64) (bool, error)

	CreateSecondKillEvent(ctx context.Context, model SecondKillEvent) error
	GetSecondKillEventByID(ctx context.Context, id int) (SecondKillEvent, error)
	ListSecondKillEvents(ctx context.Context, status string, pagination domain.Pagination) ([]SecondKillEvent, error)
	ExistsSecondKillEventByName(ctx context.Context, name string) (bool, error)
	HasUserParticipatedInSecondKill(ctx context.Context, id int, userID int64) (bool, error)

	AddParticipant(ctx context.Context, model Participant) error

	ListPendingLotteryDraws(ctx context.Context, currentTime int64) ([]LotteryDraw, error)
	UpdateLotteryDrawStatus(ctx context.Context, id int, status string) error
	ListPendingSecondKillEvents(ctx context.Context, currentTime int64) ([]SecondKillEvent, error)
	UpdateSecondKillEventStatus(ctx context.Context, id int, status string) error
	ListActiveLotteryDraws(ctx context.Context, currentTime int64) ([]LotteryDraw, error)
	ListActiveSecondKillEvents(ctx context.Context, currentTime int64) ([]SecondKillEvent, error)
}

type lotteryDrawDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

// LotteryDraw 数据库中的抽奖活动模型
type LotteryDraw struct {
	ID           int           `gorm:"primaryKey;autoIncrement"`                                                         // 抽奖活动的唯一标识符
	Name         string        `gorm:"column:name;not null"`                                                             // 抽奖活动名称
	Description  string        `gorm:"column:description;type:text"`                                                     // 抽奖活动描述
	StartTime    int64         `gorm:"column:start_time;not null"`                                                       // 活动开始时间（UNIX 时间戳）
	EndTime      int64         `gorm:"column:end_time;not null"`                                                         // 活动结束时间（UNIX 时间戳）
	Status       string        `gorm:"column:status;type:varchar(20)"`                                                   // 活动状态
	CreatedAt    int64         `gorm:"column:created_at;autoCreateTime"`                                                 // 创建时间（UNIX 时间戳）
	UpdatedAt    int64         `gorm:"column:updated_at;autoUpdateTime"`                                                 // 更新时间（UNIX 时间戳）
	Participants []Participant `gorm:"foreignKey:LotteryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 参与者列表
}

// SecondKillEvent 数据库中的秒杀活动模型
type SecondKillEvent struct {
	ID           int           `gorm:"primaryKey;autoIncrement"`                                                            // 秒杀活动的唯一标识符
	Name         string        `gorm:"column:name;not null"`                                                                // 秒杀活动名称
	Description  string        `gorm:"column:description;type:text"`                                                        // 秒杀活动描述
	StartTime    int64         `gorm:"column:start_time;not null"`                                                          // 活动开始时间（UNIX 时间戳）
	EndTime      int64         `gorm:"column:end_time;not null"`                                                            // 活动结束时间（UNIX 时间戳）
	Status       string        `gorm:"column:status;type:varchar(20)"`                                                      // 活动状态
	CreatedAt    int64         `gorm:"column:created_at;autoCreateTime"`                                                    // 创建时间（UNIX 时间戳）
	UpdatedAt    int64         `gorm:"column:updated_at;autoUpdateTime"`                                                    // 更新时间（UNIX 时间戳）
	Participants []Participant `gorm:"foreignKey:SecondKillID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 参与者列表
}

// Participant 数据库中的参与者记录模型
type Participant struct {
	ID             string `gorm:"primaryKey;column:id;type:char(36)"` // 参与记录的唯一标识符 (UUID)
	LotteryID      *int   `gorm:"column:lottery_id"`                  // 抽奖活动ID，可为null
	SecondKillID   *int   `gorm:"column:second_kill_id"`              // 秒杀活动ID，可为null
	UserID         int64  `gorm:"column:user_id;not null"`            // 参与者的用户ID
	ParticipatedAt int64  `gorm:"column:participated_at;not null"`    // 参与时间（UNIX 时间戳）
}

func NewLotteryDrawDAO(db *gorm.DB, l *zap.Logger) LotteryDrawDAO {
	return &lotteryDrawDAO{
		db: db,
		l:  l,
	}
}

// CreateLotteryDraw 创建一个新的抽奖活动
func (l *lotteryDrawDAO) CreateLotteryDraw(ctx context.Context, model LotteryDraw) error {
	if err := l.db.WithContext(ctx).Create(&model).Error; err != nil {
		l.l.Error("创建抽奖活动失败", zap.Error(err))
		return err
	}

	return nil
}

// GetLotteryDrawByID 根据ID获取指定的抽奖活动
func (l *lotteryDrawDAO) GetLotteryDrawByID(ctx context.Context, id int) (LotteryDraw, error) {
	var lotteryDraw LotteryDraw

	// 使用 Preload 预加载参与者，避免 N+1 查询问题
	if err := l.db.WithContext(ctx).
		Preload("Participants").
		Where("id = ?", id).
		First(&lotteryDraw).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.l.Warn("未找到指定ID的抽奖活动", zap.Int("ID", id))
			return LotteryDraw{}, err
		}

		l.l.Error("获取抽奖活动失败", zap.Error(err))

		return LotteryDraw{}, err
	}

	return lotteryDraw, nil
}

// ListLotteryDraws 获取所有抽奖活动，支持状态过滤和分页
func (l *lotteryDrawDAO) ListLotteryDraws(ctx context.Context, status string, pagination domain.Pagination) ([]LotteryDraw, error) {
	var lotteryDraws []LotteryDraw
	var defaultSize int64 = 10

	query := l.db.WithContext(ctx).Preload("Participants")

	// 根据状态进行过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 应用分页
	if pagination.Page <= 0 {
		pagination.Page = 1
	}

	if *pagination.Size <= 0 {
		pagination.Size = &defaultSize
	}

	query = query.Limit(int(*pagination.Size)).Offset(int(*pagination.Offset))

	if err := query.Find(&lotteryDraws).Error; err != nil {
		l.l.Error("获取抽奖活动列表失败", zap.Error(err))
		return nil, err
	}

	return lotteryDraws, nil
}

// ExistsLotteryDrawByName 检查抽奖活动名称是否存在
func (l *lotteryDrawDAO) ExistsLotteryDrawByName(ctx context.Context, name string) (bool, error) {
	var count int64

	if err := l.db.WithContext(ctx).
		Model(&LotteryDraw{}).
		Where("name = ?", name).
		Count(&count).Error; err != nil {
		l.l.Error("检查抽奖活动名称是否存在失败", zap.Error(err))
		return false, err
	}

	return count > 0, nil
}

// HasUserParticipatedInLottery 检查用户是否已参与某个抽奖活动
func (l *lotteryDrawDAO) HasUserParticipatedInLottery(ctx context.Context, id int, userID int64) (bool, error) {
	var count int64

	if err := l.db.WithContext(ctx).
		Model(&Participant{}).
		Where("lottery_id = ? AND user_id = ?", id, userID).
		Count(&count).Error; err != nil {
		l.l.Error("检查用户是否已参与抽奖活动失败", zap.Error(err))
		return false, err
	}

	return count > 0, nil
}

// CreateSecondKillEvent 创建一个新的秒杀活动
func (l *lotteryDrawDAO) CreateSecondKillEvent(ctx context.Context, model SecondKillEvent) error {
	if err := l.db.WithContext(ctx).Create(&model).Error; err != nil {
		l.l.Error("创建秒杀活动失败", zap.Error(err))
		return err
	}

	return nil
}

// GetSecondKillEventByID 根据ID获取指定的秒杀活动
func (l *lotteryDrawDAO) GetSecondKillEventByID(ctx context.Context, id int) (SecondKillEvent, error) {
	var secondKillEvent SecondKillEvent

	// 使用 Preload 预加载参与者，避免 N+1 查询问题
	if err := l.db.WithContext(ctx).
		Preload("Participants").
		First(&secondKillEvent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.l.Warn("未找到指定ID的秒杀活动", zap.Int("ID", id))
			return SecondKillEvent{}, err
		}
		l.l.Error("获取秒杀活动失败", zap.Error(err))
		return SecondKillEvent{}, err
	}

	return secondKillEvent, nil
}

// ListSecondKillEvents 获取所有秒杀活动，支持状态过滤和分页
func (l *lotteryDrawDAO) ListSecondKillEvents(ctx context.Context, status string, pagination domain.Pagination) ([]SecondKillEvent, error) {
	var secondKillEvents []SecondKillEvent
	var defaultSize int64 = 10

	query := l.db.WithContext(ctx).Preload("Participants")

	// 根据状态进行过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 应用分页
	if pagination.Page <= 0 {
		pagination.Page = 1
	}

	if *pagination.Size <= 0 {
		pagination.Size = &defaultSize
	}

	query = query.Limit(int(*pagination.Size)).Offset(int(*pagination.Offset))

	if err := query.Find(&secondKillEvents).Error; err != nil {
		l.l.Error("获取秒杀活动列表失败", zap.Error(err))
		return nil, err
	}

	return secondKillEvents, nil
}

// ExistsSecondKillEventByName 检查秒杀活动名称是否存在
func (l *lotteryDrawDAO) ExistsSecondKillEventByName(ctx context.Context, name string) (bool, error) {
	var count int64

	if err := l.db.WithContext(ctx).
		Model(&SecondKillEvent{}).
		Where("name = ?", name).
		Count(&count).Error; err != nil {
		l.l.Error("检查秒杀活动名称是否存在失败", zap.Error(err))
		return false, err
	}

	return count > 0, nil
}

// HasUserParticipatedInSecondKill 检查用户是否已参与某个秒杀活动
func (l *lotteryDrawDAO) HasUserParticipatedInSecondKill(ctx context.Context, id int, userID int64) (bool, error) {
	var count int64

	if err := l.db.WithContext(ctx).
		Model(&Participant{}).
		Where("second_kill_id = ? AND user_id = ?", id, userID).
		Count(&count).Error; err != nil {
		l.l.Error("检查用户是否已参与秒杀活动失败", zap.Error(err))
		return false, err
	}

	return count > 0, nil
}

// AddParticipant 添加参与者
func (l *lotteryDrawDAO) AddParticipant(ctx context.Context, model Participant) error {
	// 插入参与者记录
	if err := l.db.WithContext(ctx).Create(&model).Error; err != nil {
		l.l.Error("添加参与者记录失败", zap.Error(err), zap.Any("participant", model))
		return err
	}

	return nil
}

// ListPendingLotteryDraws 获取所有待激活的抽奖活动
func (l *lotteryDrawDAO) ListPendingLotteryDraws(ctx context.Context, currentTime int64) ([]LotteryDraw, error) {
	var lotteryDraws []LotteryDraw

	if err := l.db.WithContext(ctx).
		Preload("Participants").
		Where("status = ? AND start_time <= ?", domain.LotteryStatusPending, currentTime).
		Find(&lotteryDraws).Error; err != nil {
		l.l.Error("获取待激活抽奖活动失败", zap.Error(err))
		return nil, err
	}

	return lotteryDraws, nil
}

// UpdateLotteryDrawStatus 更新抽奖活动的状态
func (l *lotteryDrawDAO) UpdateLotteryDrawStatus(ctx context.Context, id int, status string) error {
	if err := l.db.WithContext(ctx).
		Model(&LotteryDraw{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		l.l.Error("更新抽奖活动状态失败", zap.Int("ID", id), zap.String("status", status), zap.Error(err))
		return err
	}

	return nil
}

// ListPendingSecondKillEvents 获取所有待激活的秒杀活动
func (l *lotteryDrawDAO) ListPendingSecondKillEvents(ctx context.Context, currentTime int64) ([]SecondKillEvent, error) {
	var secondKillEvents []SecondKillEvent

	if err := l.db.WithContext(ctx).
		Preload("Participants").
		Where("status = ? AND start_time <= ?", domain.SecondKillStatusPending, currentTime).
		Find(&secondKillEvents).Error; err != nil {
		l.l.Error("获取待激活秒杀活动失败", zap.Error(err))
		return nil, err
	}

	return secondKillEvents, nil
}

// UpdateSecondKillEventStatus 更新秒杀活动的状态
func (l *lotteryDrawDAO) UpdateSecondKillEventStatus(ctx context.Context, id int, status string) error {
	if err := l.db.WithContext(ctx).
		Model(&SecondKillEvent{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		l.l.Error("更新秒杀活动状态失败", zap.Int("ID", id), zap.String("status", status), zap.Error(err))
		return err
	}

	return nil
}

// ListActiveLotteryDraws 获取所有进行中的抽奖活动
func (l *lotteryDrawDAO) ListActiveLotteryDraws(ctx context.Context, currentTime int64) ([]LotteryDraw, error) {
	var lotteryDraws []LotteryDraw

	if err := l.db.WithContext(ctx).
		Preload("Participants").
		Where("status = ? AND start_time <= ? AND end_time >= ?", domain.LotteryStatusActive, currentTime, currentTime).
		Find(&lotteryDraws).Error; err != nil {
		l.l.Error("获取进行中的抽奖活动失败", zap.Error(err))
		return nil, err
	}

	return lotteryDraws, nil
}

// ListActiveSecondKillEvents 获取所有进行中的秒杀活动
func (l *lotteryDrawDAO) ListActiveSecondKillEvents(ctx context.Context, currentTime int64) ([]SecondKillEvent, error) {
	var secondKillEvents []SecondKillEvent

	if err := l.db.WithContext(ctx).
		Preload("Participants").
		Where("status = ? AND start_time <= ? AND end_time >= ?", domain.SecondKillStatusActive, currentTime, currentTime).
		Find(&secondKillEvents).Error; err != nil {
		l.l.Error("获取进行中的秒杀活动失败", zap.Error(err))
		return nil, err
	}

	return secondKillEvents, nil
}
