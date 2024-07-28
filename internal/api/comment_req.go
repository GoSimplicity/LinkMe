package api

type CreateCommentReq struct {
	PostId  int64  `json:"postId" binding:"required"`
	Content string `json:"content" binding:"required"`
	RootId  *int64 `json:"root_id,omitempty"` // 根评论ID，顶层评论时为空
	PID     *int64 `json:"pid,omitempty"`     // 父评论ID，顶层评论时为空
}

type ListCommentsReq struct {
	PostId int64
	MinId  int64
	Limit  int64
}

type DeleteCommentReq struct {
	CommentId int64 `json:"commentId" binding:"required"`
}

type GetMoreCommentReplyReq struct {
	RootId int64
	MaxId  int64
	Limit  int64
}
