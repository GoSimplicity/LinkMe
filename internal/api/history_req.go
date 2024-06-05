package api

type ListHistoryReq struct {
	Page int    `json:"page,omitempty"` // 当前页码
	Size *int64 `json:"size,omitempty"` // 每页数据量
}

type DeleteHistoryReq struct {
	ID int64 `uri:"historyId"`
}
