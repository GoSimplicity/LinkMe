package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/req"
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
	permissionGroup.POST("/search_comment", WrapBody(s.SearchComment))
}

func (s *SearchHandler) SearchUser(ctx *gin.Context, req req.SearchReq) (Result, error) {
	users, err := s.svc.SearchUsers(ctx, req.Expression)
	if err != nil {
		return Result{
			Code: SearchUserERRORCode,
			Msg:  SearchUserERROR,
		}, nil
	}
	return Result{
		Code: RequestsOK,
		Msg:  SearchUserSuccess,
		Data: users,
	}, nil
}

func (s *SearchHandler) SearchPost(ctx *gin.Context, req req.SearchReq) (Result, error) {
	posts, err := s.svc.SearchPosts(ctx, req.Expression)
	if err != nil {
		return Result{
			Code: SearchPostERRORCode,
			Msg:  SearchPostERROR,
		}, nil
	}
	return Result{
		Code: RequestsOK,
		Msg:  SearchPostSuccess,
		Data: posts,
	}, nil
}
func (s *SearchHandler) SearchComment(ctx *gin.Context, req req.SearchReq) (Result, error) {
	comments, err := s.svc.SearchComments(ctx, req.Expression)
	if err != nil {
		return Result{
			Code: SearchCommentERRORCode,
			Msg:  SearchCommentERROR,
		}, nil
	}
	return Result{
		Code: RequestsOK,
		Msg:  SearchCommentSuccess,
		Data: comments,
	}, nil
}
