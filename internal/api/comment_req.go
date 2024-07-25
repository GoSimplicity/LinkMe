package api

type CreateCommentReq struct {
	Content string `json:"content" binding:"required"`
	PostId  int64  `json:"postId" binding:"required"`
}

type ListCommentsReq struct {
}

type DeleteCommentReq struct {
}

type GetMoreCommentReplyReq struct {
}
