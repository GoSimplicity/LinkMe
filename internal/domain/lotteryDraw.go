package domain

import "errors"

var ErrNotFound = errors.New("record not found")

const (
	LotteryStatusPending   string = "pending"   // 待开始
	LotteryStatusActive    string = "active"    // 进行中
	LotteryStatusCompleted string = "completed" // 已完成
)

const (
	SecondKillStatusPending   string = "pending"   // 待开始
	SecondKillStatusActive    string = "active"    // 进行中
	SecondKillStatusCompleted string = "completed" // 已完成
)

// Participant 表示参与者的记录，适用于抽奖和秒杀活动
type Participant struct {
	ID             string // 参与记录的唯一标识符
	LotteryID      *int   // 关联的活动ID（可以是抽奖或秒杀活动）
	SecondKillID   *int
	ActivityType   string
	UserID         int64 // 参与者的用户ID
	ParticipatedAt int64 // UNIX 时间戳，表示参与时间
}

// LotteryDraw 表示一个抽奖活动
type LotteryDraw struct {
	ID           int           // 抽奖活动的唯一标识符
	Name         string        // 抽奖活动名称
	Description  string        // 抽奖活动描述
	StartTime    int64         // UNIX 时间戳，表示活动开始时间
	EndTime      int64         // UNIX 时间戳，表示活动结束时间
	Status       string        // 抽奖活动状态
	Participants []Participant // 参与者列表
}

// SecondKillEvent 表示一个秒杀活动
type SecondKillEvent struct {
	ID           int           // 秒杀活动的唯一标识符
	Name         string        // 秒杀活动名称
	Description  string        // 秒杀活动描述
	StartTime    int64         // UNIX 时间戳，表示活动开始时间
	EndTime      int64         // UNIX 时间戳，表示活动结束时间
	Status       string        // 秒杀活动状态
	Participants []Participant // 参与者列表
}
