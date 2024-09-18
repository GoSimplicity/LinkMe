package dao

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LotteryDrawDAO interface {
	// 抽奖相关的方法

	// CreateLotteryDraw 创建一个新的抽奖活动
	CreateLotteryDraw(model LotteryDrawModel) (LotteryDrawModel, error)

	// GetLotteryDrawByID 根据ID获取指定的抽奖活动
	GetLotteryDrawByID(id string) (LotteryDrawModel, error)

	// GetAllLotteryDraws 获取所有抽奖活动
	GetAllLotteryDraws() ([]LotteryDrawModel, error)

	// 秒杀相关的方法

	// CreateSecondKillEvent 创建一个新的秒杀活动
	CreateSecondKillEvent(model SecondKillEventModel) (SecondKillEventModel, error)

	// GetSecondKillEventByID 根据ID获取指定的秒杀活动
	GetSecondKillEventByID(id string) (SecondKillEventModel, error)

	// GetAllSecondKillEvents 获取所有秒杀活动
	GetAllSecondKillEvents() ([]SecondKillEventModel, error)

	// 参与者相关的方法

	// AddParticipant 添加一个参与者记录
	AddParticipant(model ParticipantModel) (ParticipantModel, error)
}

type lotteryDrawDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

// LotteryDrawModel 表示数据库中的抽奖活动模型
type LotteryDrawModel struct {
	ID           string             `gorm:"primaryKey;autoIncrement"`                                            // 抽奖活动的唯一标识符
	Name         string             `gorm:"column:name;not null"`                                                // 抽奖活动名称
	Description  string             `gorm:"column:description;type:text"`                                        // 抽奖活动描述
	StartTime    int64              `gorm:"column:start_time;not null"`                                          // 活动开始时间（UNIX 时间戳）
	EndTime      int64              `gorm:"column:end_time;not null"`                                            // 活动结束时间（UNIX 时间戳）
	Status       string             `gorm:"column:status;type:varchar(20)"`                                      // 活动状态
	CreatedAt    int64              `gorm:"column:created_at;autoCreateTime"`                                    // 创建时间（UNIX 时间戳）
	UpdatedAt    int64              `gorm:"column:updated_at;autoUpdateTime"`                                    // 更新时间（UNIX 时间戳）
	Participants []ParticipantModel `gorm:"foreignKey:ActivityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 参与者列表
}

// SecondKillEventModel 表示数据库中的秒杀活动模型
type SecondKillEventModel struct {
	ID           string             `gorm:"primaryKey;autoIncrement"`                                            // 秒杀活动的唯一标识符
	Name         string             `gorm:"column:name;not null"`                                                // 秒杀活动名称
	Description  string             `gorm:"column:description;type:text"`                                        // 秒杀活动描述
	StartTime    int64              `gorm:"column:start_time;not null"`                                          // 活动开始时间（UNIX 时间戳）
	EndTime      int64              `gorm:"column:end_time;not null"`                                            // 活动结束时间（UNIX 时间戳）
	Status       string             `gorm:"column:status;type:varchar(20)"`                                      // 活动状态
	CreatedAt    int64              `gorm:"column:created_at;autoCreateTime"`                                    // 创建时间（UNIX 时间戳）
	UpdatedAt    int64              `gorm:"column:updated_at;autoUpdateTime"`                                    // 更新时间（UNIX 时间戳）
	Participants []ParticipantModel `gorm:"foreignKey:ActivityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 参与者列表
}

// ParticipantModel 表示数据库中的参与者记录模型
type ParticipantModel struct {
	ID             string `gorm:"primaryKey;autoIncrement"`        // 参与记录的唯一标识符
	ActivityID     string `gorm:"column:activity_id;not null"`     // 关联的活动ID（抽奖或秒杀）
	UserID         string `gorm:"column:user_id;not null"`         // 参与者的用户ID
	ParticipatedAt int64  `gorm:"column:participated_at;not null"` // 参与时间（UNIX 时间戳）
}

func NewLotteryDrawDAO(db *gorm.DB, l *zap.Logger) LotteryDrawDAO {
	return &lotteryDrawDAO{
		db: db,
		l:  l,
	}
}

func (l *lotteryDrawDAO) CreateLotteryDraw(model LotteryDrawModel) (LotteryDrawModel, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawDAO) GetLotteryDrawByID(id string) (LotteryDrawModel, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawDAO) GetAllLotteryDraws() ([]LotteryDrawModel, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawDAO) CreateSecondKillEvent(model SecondKillEventModel) (SecondKillEventModel, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawDAO) GetSecondKillEventByID(id string) (SecondKillEventModel, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawDAO) GetAllSecondKillEvents() ([]SecondKillEventModel, error) {
	//TODO implement me
	panic("implement me")
}

func (l *lotteryDrawDAO) AddParticipant(model ParticipantModel) (ParticipantModel, error) {
	//TODO implement me
	panic("implement me")
}
