package required_parameter

type ListHistoryReq struct {
	Page int    `json:"page,omitempty"` // 当前页码
	Size *int64 `json:"size,omitempty"` // 每页数据量
}

type DeleteHistoryReq struct {
	PostId int64 `json:"postId,omitempty"`
}

type DeleteHistoryAllReq struct {
	IsDeleteAll bool `json:"isDeleteAll,omitempty"`
}
