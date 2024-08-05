package req

type CreatePlateReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeletePlateReq struct {
	PlateID int64 `uri:"plateId"`
}

type UpdatePlateReq struct {
	ID          int64  `json:"plateId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type ListPlateReq struct {
	Page int    `json:"page,omitempty"` // 当前页码
	Size *int64 `json:"size,omitempty"` // 每页数据量
}
