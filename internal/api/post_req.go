package api

type EditReq struct {
	PostId  int64  `json:"postId,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

type PublishReq struct {
	PostId int64 `uri:"postId"`
}

type WithDrawReq struct {
	PostId int64 `uri:"postId"`
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
	Like   bool  `json:"like,omitempty"`
}

type CollectReq struct {
	PostId    int64 `json:"postId,omitempty"`
	CollectId bool  `json:"collectId,omitempty"`
}
