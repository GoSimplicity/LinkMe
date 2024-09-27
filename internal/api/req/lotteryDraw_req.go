package req

// ListLotteryDrawsReq 定义获取所有抽奖活动的请求参数
type ListLotteryDrawsReq struct {
	Page   int    `json:"page,omitempty"` // 当前页码
	Size   *int64 `json:"size,omitempty"` // 每页数据量
	Status string `json:"status"`         // 抽奖活动状态过滤
}

// CreateLotteryDrawReq 定义创建新的抽奖活动的请求参数
type CreateLotteryDrawReq struct {
	Name        string `json:"name"`        // 抽奖活动名称
	Description string `json:"description"` // 抽奖活动描述
	StartTime   int64  `json:"startTime"`   // 活动开始时间，必须晚于当前时间
	EndTime     int64  `json:"endTime"`     // 活动结束时间，必须晚于开始时间
}

// GetLotteryDrawReq 定义获取指定ID抽奖活动的请求参数
type GetLotteryDrawReq struct {
	ID int `uri:"id"` // 抽奖活动的唯一标识符
}

// ParticipateLotteryDrawReq 定义参与抽奖活动的请求参数
type ParticipateLotteryDrawReq struct {
	ActivityID   int    `json:"activityId"` // 抽奖活动的唯一标识符
	ActivityType string `json:"activityType"`
}

// GetAllSecondKillEventsReq 定义获取所有秒杀活动的请求参数
type GetAllSecondKillEventsReq struct {
	Page   int    `json:"page,omitempty"` // 当前页码
	Size   *int64 `json:"size,omitempty"` // 每页数据量
	Status string `json:"status"`         // 秒杀活动状态过滤
}

// CreateSecondKillEventReq 定义创建新的秒杀活动的请求参数
type CreateSecondKillEventReq struct {
	Name        string `json:"name"`        // 秒杀活动名称
	Description string `json:"description"` // 秒杀活动描述
	StartTime   int64  `json:"startTime"`   // 活动开始时间，必须晚于当前时间
	EndTime     int64  `json:"endTime"`     // 活动结束时间，必须晚于开始时间
}

// GetSecondKillEventReq 定义获取指定ID秒杀活动的请求参数
type GetSecondKillEventReq struct {
	ID int `uri:"id"` // 秒杀活动的唯一标识符
}

// ParticipateSecondKillReq 定义参与秒杀活动的请求参数
type ParticipateSecondKillReq struct {
	ActivityID int `json:"activityId"` // 秒杀活动的唯一标识符
}
