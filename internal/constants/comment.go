package constants

const (
	// 错误代码
	CreateCommentErrorCode       = 406001
	DeleteCommentErrorCode       = 406002
	ListCommentErrorCode         = 406003
	GetMoreCommentReplyErrorCode = 406004
	GetTopCommentReplyErrorCode  = 406005

	// 错误信息
	CreateCommentErrorMsg       = "Failed to create comment"
	DeleteCommentErrorMsg       = "Failed to delete comment"
	ListCommentErrorMsg         = "Failed to list comments"
	GetMoreCommentReplyErrorMsg = "Failed to get more comment replies"
	GetTopCommentReplyErrorMsg  = "Failed to get top comment replies"

	// 成功信息
	CreateCommentSuccessMsg       = "Comment created successfully"
	DeleteCommentSuccessMsg       = "Comment deleted successfully"
	ListCommentSuccessMsg         = "Comments listed successfully"
	GetMoreCommentReplySuccessMsg = "More comment replies retrieved successfully"
	GetTopCommentReplySuccessMsg  = "Top comment replies retrieved successfully"
)
