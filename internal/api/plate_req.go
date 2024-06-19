package api

type CreatePlateReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeletePlateReq struct {
	ID int64 `uri:"id"`
}

type UpdatePlateReq struct {
	ID          int64  `uri:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type ListPlateReq struct {
	Page int    `json:"page,omitempty"` // 当前页码
	Size *int64 `json:"size,omitempty"` // 每页数据量
}
