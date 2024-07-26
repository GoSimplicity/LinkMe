package api

type CreateCommentReq struct {
	Content string `json:"content" binding:"required"`
	PostId  int64  `json:"postId" binding:"required"`
}

type ListCommentsReq struct {
	biz    string
	bizId  int64
	min_id int64
	limit  int64
}

type DeleteCommentReq struct {
	CommentId int64 `json:"commentId" binding:"required"`
}

type GetMoreCommentReplyReq struct {
}
