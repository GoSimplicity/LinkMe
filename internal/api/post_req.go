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
