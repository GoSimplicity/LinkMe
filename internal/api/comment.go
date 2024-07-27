package api

import (
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	. "github.com/GoSimplicity/LinkMe/pkg/ginp"

	"github.com/gin-gonic/gin"
)

// 评论处理器结构体
type CommentHandler struct {
	svc service.CommentService
}

// 创建新的评论处理器
func NewCommentHandler(svc service.CommentService) *CommentHandler {
	return &CommentHandler{
		svc: svc,
	}
}

// 注册路由
func (ch *CommentHandler) RegisterRoutes(server *gin.Engine) {
	commentsGroup := server.Group("/api/comments")
	commentsGroup.POST("/create", WrapBody(ch.CreateComment))
	commentsGroup.POST("/list", WrapBody(ch.ListComments))
	commentsGroup.DELETE("/delete", WrapBody(ch.DeleteComment))
	commentsGroup.POST("/get_more", WrapBody(ch.GetMoreCommentReply))
}

// 创建评论处理器方法
func (ch *CommentHandler) CreateComment(ctx *gin.Context, req CreateCommentReq) (Result, error) {
	err := ch.svc.CreateComment(ctx, domain.Comment{
		Content: req.Content,
		PostId:  req.PostId,
	})
	if err != nil {
		return Result{
			Code: CreateCommentErrorCode,
			Msg:  CreateCommentErrorMsg,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  CreateCommentSuccessMsg,
	}, nil
}

// 列出评论处理器方法
func (ch *CommentHandler) ListComments(ctx *gin.Context, req ListCommentsReq) (Result, error) {
	comments, err := ch.svc.ListComments(ctx, req.biz, req.bizId, req.min_id, req.limit)
	if err != nil {
		return Result{
			Code: ListCommentErrorCode,
			Msg:  ListCommentErrorMsg,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  ListCommentSuccessMsg,
		Data: comments,
	}, nil
}

// 删除评论处理器方法
func (ch *CommentHandler) DeleteComment(ctx *gin.Context, req DeleteCommentReq) (Result, error) {
	err := ch.svc.DeleteComment(ctx, req.CommentId)
	if err != nil {
		return Result{
			Code: DeleteCommentErrorCode,
			Msg:  DeleteCommentErrorMsg,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  DeleteCommentSuccessMsg,
	}, nil
}

// 获取更多评论回复处理器方法
func (ch *CommentHandler) GetMoreCommentReply(ctx *gin.Context, req GetMoreCommentReplyReq) (Result, error) {
	// 由于方法未实现，返回空结果
	return Result{}, nil
}
