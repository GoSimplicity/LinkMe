package required_parameter

type CreateCommentReq struct {
	PostId  int64  `json:"postId" binding:"required"`
	Content string `json:"content" binding:"required"`
	RootId  *int64 `json:"rootId,omitempty"` // 根评论ID，顶层评论时为空
	PID     *int64 `json:"pid,omitempty"`    // 父评论ID，顶层评论时为空
}

type ListCommentsReq struct {
	PostId int64 `json:"postId"`
	MinId  int64 `json:"minId"`
	Limit  int64 `json:"limit"`
}

type DeleteCommentReq struct {
	CommentId int64 `uri:"commentId"`
}

type GetMoreCommentReplyReq struct {
	RootId int64 `json:"rootId"`
	MaxId  int64 `json:"maxId"`
	Limit  int64 `json:"limit"`
}
