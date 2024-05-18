package api

import (
	. "LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/service"
	. "LinkMe/pkg/ginp"
	ijwt "LinkMe/utils/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PostHandler struct {
	svc service.PostService
	l   *zap.Logger
}

func NewPostHandler(svc service.PostService, l *zap.Logger) *PostHandler {
	return &PostHandler{
		svc: svc,
		l:   l,
	}
}
func (ph *PostHandler) RegisterRoutes(server *gin.Engine) {
	postGroup := server.Group("/posts")
	postGroup.POST("/edit", WrapBody(ph.Edit))         // 编辑帖子
	postGroup.PUT("/update", WrapBody(ph.Update))      // 更新帖子
	postGroup.PUT("/publish", WrapBody(ph.Publish))    // 更新帖子状态为发布
	postGroup.PUT("/withdraw", WrapBody(ph.Withdraw))  // 更新帖子状态为撤回
	postGroup.GET("/list", WrapBody(ph.List))          // 可以添加分页和排序参数
	postGroup.GET("/list_pub", WrapBody(ph.ListPub))   // 同上
	postGroup.GET("/detail/:postId", ph.Detail)        // 使用参数获取特定帖子
	postGroup.GET("/detail_pub/:postId", ph.DetailPub) // 同上
	postGroup.DELETE("/:postId", ph.DeletePost)        // 使用 DELETE 方法删除特定帖子
	postGroup.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "hello")
	})
}

func (ph *PostHandler) Edit(ctx *gin.Context, req EditReq) (Result, error) {
	// 获取当前登陆的用户信息
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	id, err := ph.svc.Create(ctx, domain.Post{
		ID:      req.PostId,
		Content: req.Content,
		Title:   req.Title,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ph.l.Error(PostEditERROR, zap.Error(err))
		return Result{
			Code: PostInternalServerError,
			Msg:  PostServerERROR,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  PostEditSuccess,
		Data: id,
	}, nil
}

func (ph *PostHandler) Update(ctx *gin.Context, req UpdateReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if err := ph.svc.Update(ctx, domain.Post{
		ID: req.PostId,
		Author: domain.Author{
			Id: uc.Uid,
		},
	}); err != nil {
		ph.l.Error(PostUpdateERROR, zap.Error(err))
		return Result{
			Code: PostInternalServerError,
			Msg:  PostServerERROR,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  PostUpdateSuccess,
	}, nil
}

func (ph *PostHandler) Publish(ctx *gin.Context, req PublishReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if err := ph.svc.Publish(ctx, req.PostId, domain.Post{
		ID: req.PostId,
		Author: domain.Author{
			Id: uc.Uid,
		},
	}); err != nil {
		ph.l.Error(PostPublishERROR, zap.Error(err))
		return Result{
			Code: PostInternalServerError,
			Msg:  PostServerERROR,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  PostPublishSuccess,
		Data: req.PostId,
	}, nil
}

func (ph *PostHandler) Withdraw(ctx *gin.Context, req WithDrawReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if err := ph.svc.Withdraw(ctx, req.PostId, domain.Post{
		ID: req.PostId,
		Author: domain.Author{
			Id: uc.Uid,
		},
	}); err != nil {
		ph.l.Error(PostWithdrawERROR, zap.Error(err))
		return Result{
			Code: PostInternalServerError,
			Msg:  PostServerERROR,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  PostWithdrawSuccess,
		Data: req.PostId,
	}, nil
}

func (ph *PostHandler) List(ctx *gin.Context, req ListReq) (Result, error) {
	return Result{}, nil
}

func (ph *PostHandler) ListPub(ctx *gin.Context, req ListPubReq) (Result, error) {
	du, err := ph.svc.ListPublishedPosts(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		ph.l.Error(PostListPubERROR, zap.Error(err))
		return Result{
			Code: PostInternalServerError,
			Msg:  PostServerERROR,
		}, err
	}
	return Result{
		Code: RequestsOK,
		Msg:  PostListPubSuccess,
		Data: du,
	}, nil
}

func (ph *PostHandler) Detail(ctx *gin.Context) {

}

func (ph *PostHandler) DetailPub(ctx *gin.Context) {

}

func (ph *PostHandler) DeletePost(ctx *gin.Context) {

}
