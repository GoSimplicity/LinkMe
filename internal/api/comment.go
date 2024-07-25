package api

import (
	"LinkMe/internal/service"
	. "LinkMe/pkg/ginp"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	svc service.CommentService
}

func NewCommentHandler(svc service.CommentService) *CommentHandler {
	return &CommentHandler{
		svc: svc,
	}
}

func (ch *CommentHandler) RegisterRoutes(server *gin.Engine) {
	historyGroup := server.Group("/api/comments")
	historyGroup.POST("/create", WrapBody(ch.CreateComment))
	historyGroup.POST("/list", WrapBody(ch.ListComments))
	historyGroup.DELETE("/delete", WrapBody(ch.DeleteComment))
	historyGroup.POST("/get_more", WrapBody(ch.GetMoreCommentReply))
}

func (ch *CommentHandler) CreateComment(ctx *gin.Context, req CreateCommentReq) (Result, error) {
	return Result{}, nil
}

func (ch *CommentHandler) ListComments(ctx *gin.Context, req ListCommentsReq) (Result, error) {
	return Result{}, nil
}

func (ch *CommentHandler) DeleteComment(ctx *gin.Context, req DeleteCommentReq) (Result, error) {
	return Result{}, nil
}

func (ch *CommentHandler) GetMoreCommentReply(ctx *gin.Context, req GetMoreCommentReplyReq) (Result, error) {
	return Result{}, nil
}
