package api

type EditReq struct {
	PostId  int64  `json:"postId,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

type PublishReq struct {
	PostId int64 `json:"postId,omitempty"`
}

type WithDrawReq struct {
	PostId int64 `json:"postId,omitempty"`
}

type ListReq struct {
	Page int    `json:"page,omitempty"` // 当前页码
	Size *int64 `json:"size,omitempty"` // 每页数据量
}

type ListPubReq struct {
	Page int    `json:"page,omitempty"` // 当前页码
	Size *int64 `json:"size,omitempty"` // 每页数据量
}

type UpdateReq struct {
	PostId  int64  `json:"postId,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}
type DetailReq struct {
	PostId int64 `uri:"postId"`
}

type DeleteReq struct {
	PostId int64 `uri:"postId"`
}

type LikeReq struct {
	PostId int64 `json:"postId,omitempty"`
	Liked  bool  `json:"liked,omitempty"`
}

type CollectReq struct {
	PostId    int64 `json:"postId,omitempty"`
	CollectId int64 `json:"collectId,omitempty"`
	Collectd  bool  `json:"collectd,omitempty"`
}

type InteractReq struct {
	BizId   []int64 `json:"bizId,omitempty"`
	BizName string  `json:"bizName,omitempty"`
}
