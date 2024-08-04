package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/required_parameter"
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	. "github.com/GoSimplicity/LinkMe/pkg/ginp"
	"github.com/gin-gonic/gin"
)

type RelationHandler struct {
	svc service.RelationService
}

func (r *RelationHandler) RegisterRoutes(server *gin.Engine) {
	relationGroup := server.Group("/api/relations")
	relationGroup.POST("/list", WrapBody(r.ListRelations))                  // 查看用户关系列表
	relationGroup.GET("/get_info", WrapQuery(r.GetRelationInfo))            // 查看用户关系信息
	relationGroup.GET("/get_followee_count", WrapQuery(r.GetFolloweeCount)) // 获取关注者数量
	relationGroup.GET("/get_follower_count", WrapQuery(r.GetFollowerCount)) // 获取粉丝数量
	relationGroup.POST("/follow", WrapBody(r.FollowUser))                   // 关注
	relationGroup.POST("/cancel_follow", WrapBody(r.CancelFollowUser))      // 关注
}

func NewRelationHandler(svc service.RelationService) *RelationHandler {
	return &RelationHandler{
		svc: svc,
	}
}

// ListRelations 处理列出关注关系的请求
func (r *RelationHandler) ListRelations(ctx *gin.Context, req required_parameter.ListRelationsReq) (Result, error) {
	relations, err := r.svc.ListRelations(ctx, req.FollowerID, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return Result{
			Code: ListCommentErrorCode,
			Msg:  ListCommentErrorMsg,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  ListCommentSuccessMsg,
		Data: relations,
	}, nil
}

// GetRelationInfo 处理获取关注关系信息的请求
func (r *RelationHandler) GetRelationInfo(ctx *gin.Context, req required_parameter.GetRelationInfoReq) (Result, error) {
	relation, err := r.svc.GetRelationInfo(ctx, req.FollowerID, req.FolloweeID)
	if err != nil {
		return Result{
			Code: ListCommentErrorCode,
			Msg:  ListCommentErrorMsg,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  ListCommentSuccessMsg,
		Data: relation,
	}, nil
}

// FollowUser 处理关注用户的请求
func (r *RelationHandler) FollowUser(ctx *gin.Context, req required_parameter.FollowUserReq) (Result, error) {
	if err := r.svc.FollowUser(ctx, req.FollowerID, req.FolloweeID); err != nil {
		return Result{
			Code: FollowUserERRORCode,
			Msg:  FollowUserERRORMsg,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  FollowUserSuccessMsg,
	}, nil
}

func (r *RelationHandler) CancelFollowUser(ctx *gin.Context, req required_parameter.CancelFollowUserReq) (Result, error) {
	if err := r.svc.CancelFollowUser(ctx, req.FollowerID, req.FolloweeID); err != nil {
		return Result{
			Code: CancelFollowUserERRORCode,
			Msg:  CancelFollowUserERRORMsg,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  CancelFollowUserSuccessMsg,
	}, nil
}

// GetFolloweeCount 获取关注者的数量
func (r *RelationHandler) GetFolloweeCount(ctx *gin.Context, req required_parameter.GetFolloweeCountReq) (Result, error) {
	count, err := r.svc.GetFolloweeCount(ctx, req.UserID)
	if err != nil {
		return Result{
			Code: GetFolloweeCountErrorCode,
			Msg:  GetFolloweeCountERRORMsg,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  GetFolloweeCountSuccessMsg,
		Data: count,
	}, nil
}

// GetFollowerCount 获取粉丝的数量
func (r *RelationHandler) GetFollowerCount(ctx *gin.Context, req required_parameter.GetFollowerCountReq) (Result, error) {
	count, err := r.svc.GetFollowerCount(ctx, req.UserID)
	if err != nil {
		return Result{
			Code: GetFollowerCountErrorCode,
			Msg:  GetFollowerCountERRORMsg,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  GetFollowerCountSuccessMsg,
		Data: count,
	}, nil
}
