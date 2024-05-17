package api

type EditReq struct {
	Id      int64  `json:"id,omitempty"`
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
