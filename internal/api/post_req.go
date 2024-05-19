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
