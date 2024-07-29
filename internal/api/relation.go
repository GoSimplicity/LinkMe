package api

import (
	"github.com/GoSimplicity/LinkMe/internal/service"
	. "github.com/GoSimplicity/LinkMe/pkg/ginp"
	"github.com/gin-gonic/gin"
)

type RelationHandler struct {
	svc service.RelationService
}

func (r *RelationHandler) RegisterRoutes(server *gin.Engine) {
	relationGroup := server.Group("/api/relations")
	relationGroup.POST("/list", WrapBody(r.ListRelations))
	relationGroup.GET("/get_info", WrapQuery(r.GetRelationInfo))
	relationGroup.POST("/follow", WrapBody(r.FollowUser))
}

func NewRelationHandler(svc service.RelationService) *RelationHandler {
	return &RelationHandler{
		svc: svc,
	}
}

// ListRelations 处理列出关注关系的请求
func (r *RelationHandler) ListRelations(ctx *gin.Context, req ListRelationsReq) (Result, error) {
	// TODO 实现方法
	return Result{}, nil
}

// GetRelationInfo 处理获取关注关系信息的请求
func (r *RelationHandler) GetRelationInfo(ctx *gin.Context, req GetRelationInfoReq) (Result, error) {
	// TODO 实现方法
	return Result{}, nil
}

// FollowUser 处理关注用户的请求
func (r *RelationHandler) FollowUser(ctx *gin.Context, req FollowUserReq) (Result, error) {
	// TODO 实现方法
	return Result{}, nil
}
