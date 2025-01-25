package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/req"
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	. "github.com/GoSimplicity/LinkMe/pkg/ginp"
	icontentfilter "github.com/GoSimplicity/LinkMe/utils/contentfilter"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/gin-gonic/gin"
)

// CommentHandler 评论处理器结构体
type CommentHandler struct {
	svc service.CommentService
}

// NewCommentHandler 创建新的评论处理器
func NewCommentHandler(svc service.CommentService) *CommentHandler {
	return &CommentHandler{
		svc: svc,
	}
}

// RegisterRoutes 注册路由
func (ch *CommentHandler) RegisterRoutes(server *gin.Engine) {
	commentsGroup := server.Group("/api/comments")
	commentsGroup.POST("/create", WrapBody(ch.CreateComment))
	commentsGroup.POST("/list", WrapBody(ch.ListComments))
	commentsGroup.DELETE("/delete/:commentId", WrapParam(ch.DeleteComment))
	commentsGroup.POST("/get_more", WrapBody(ch.GetMoreCommentReply))
}

// CreateComment 创建评论处理器方法
func (ch *CommentHandler) CreateComment(ctx *gin.Context, req req.CreateCommentReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	// 进行敏感词过滤
	SensitiveContent := icontentfilter.SensitiveFilterFun(req.Content)
	comment := domain.Comment{
		Content: SensitiveContent,
		PostId:  req.PostId,
		UserId:  uc.Uid,
		Biz:     "comment",
		BizId:   req.PostId,
	}
	if req.RootId != nil {
		comment.RootComment = &domain.Comment{Id: *req.RootId}
	}
	if req.PID != nil {
		comment.ParentComment = &domain.Comment{Id: *req.PID}
	}

	err := ch.svc.CreateComment(ctx, comment)
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

// ListComments 列出评论处理器方法
func (ch *CommentHandler) ListComments(ctx *gin.Context, req req.ListCommentsReq) (Result, error) {
	comments, err := ch.svc.ListComments(ctx, req.PostId, req.MinId, req.Limit)
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

// DeleteComment 删除评论处理器方法
func (ch *CommentHandler) DeleteComment(ctx *gin.Context, req req.DeleteCommentReq) (Result, error) {
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

// GetMoreCommentReply 获取更多评论回复处理器方法
func (ch *CommentHandler) GetMoreCommentReply(ctx *gin.Context, req req.GetMoreCommentReplyReq) (Result, error) {
	comments, err := ch.svc.GetMoreCommentsReply(ctx, req.RootId, req.MaxId, req.Limit)
	if err != nil {
		return Result{
			Code: GetMoreCommentReplyErrorCode,
			Msg:  GetMoreCommentReplyErrorMsg,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  GetMoreCommentReplySuccessMsg,
		Data: comments,
	}, nil
}
