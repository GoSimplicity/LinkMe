package api

import (
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/service"
	. "github.com/GoSimplicity/LinkMe/pkg/ginp"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	svc service.SearchService
}

func NewSearchHandler(svc service.SearchService) *SearchHandler {
	return &SearchHandler{
		svc: svc,
	}
}

func (s *SearchHandler) RegisterRoutes(server *gin.Engine) {
	permissionGroup := server.Group("/api/search")
	permissionGroup.POST("/search_user", WrapBody(s.SearchUser))
	permissionGroup.POST("/search_post", WrapBody(s.SearchPost))
}

func (s *SearchHandler) SearchUser(ctx *gin.Context, req SearchReq) (Result, error) {
	users, err := s.svc.SearchUsers(ctx, req.userID, req.expression)
	if err != nil {
		return Result{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}
	return Result{
		Code: RequestsOK,
		Msg:  "Success",
		Data: users,
	}, nil
}

func (s *SearchHandler) SearchPost(ctx *gin.Context, req SearchReq) (Result, error) {
	posts, err := s.svc.SearchPosts(ctx, req.userID, req.expression)
	if err != nil {
		return Result{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}
	return Result{
		Code: RequestsOK,
		Msg:  "Success",
		Data: posts,
	}, nil
}
