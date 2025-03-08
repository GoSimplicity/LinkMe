package req

type EditReq struct {
	PostId  uint   `json:"postId,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	PlateID int64  `json:"plateId,omitempty"`
}

type PublishReq struct {
	PostId uint `json:"postId,omitempty"`
}

type WithDrawReq struct {
	PostId uint `json:"postId,omitempty"`
}

type ListReq struct {
	Page int    `json:"page,omitempty"` // 当前页码
	Size *int64 `json:"size,omitempty"` // 每页数据量
}

type DetailPostReq struct {
	PostId uint `uri:"postId"`
}

type UpdateReq struct {
	PostId  uint   `json:"postId,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	PlateID int64  `json:"plateId,omitempty"`
}
type DetailReq struct {
	PostId uint `uri:"postId"`
}

type DeleteReq struct {
	PostId uint `uri:"postId"`
}

type LikeReq struct {
	PostId uint `json:"postId,omitempty"`
	Liked  bool `json:"liked,omitempty"`
}

type CollectReq struct {
	PostId   uint `json:"postId,omitempty"`
	Collectd bool `json:"collectd,omitempty"`
}

// type InteractReq struct {
// 	BizId   []int64 `json:"bizId,omitempty"`
// 	BizName string  `json:"bizName,omitempty"`
// }

// type GetPostCountReq struct {
// }

type SearchByPlateReq struct {
	PlateId int64  `json:"plateId,omitempty"`
	Page    int    `json:"page,omitempty"`
	Size    *int64 `json:"size,omitempty"`
}
