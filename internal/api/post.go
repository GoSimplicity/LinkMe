package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/req"
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/middleware"
	. "github.com/GoSimplicity/LinkMe/pkg/ginp"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	svc    service.PostService
	intSvc service.InteractiveService
	ce     *casbin.Enforcer
	biz    string
}

func NewPostHandler(svc service.PostService, intSvc service.InteractiveService, ce *casbin.Enforcer) *PostHandler {
	return &PostHandler{
		svc:    svc,
		intSvc: intSvc,
		ce:     ce,
		biz:    "post",
	}
}

func (ph *PostHandler) RegisterRoutes(server *gin.Engine) {
	casbinMiddleware := middleware.NewCasbinMiddleware(ph.ce)
	postGroup := server.Group("/api/posts")

	// 基础帖子操作
	postGroup.POST("/edit", WrapBody(ph.Edit))                    // 编辑帖子
	postGroup.POST("/update", WrapBody(ph.Update))                // 更新帖子
	postGroup.POST("/publish", WrapBody(ph.Publish))              // 发布帖子
	postGroup.POST("/withdraw", WrapBody(ph.Withdraw))            // 撤回帖子
	postGroup.DELETE("/delete/:postId", WrapParam(ph.DeletePost)) // 删除帖子

	// 帖子列表相关
	postGroup.POST("/list", WrapBody(ph.List))                                          // 获取个人帖子列表
	postGroup.POST("/list_pub", WrapBody(ph.ListPub))                                   // 获取公开帖子列表
	postGroup.POST("/list_post", casbinMiddleware.CheckCasbin(), WrapBody(ph.ListPost)) // 管理员获取所有帖子

	// 帖子详情相关
	postGroup.GET("/detail/:postId", WrapParam(ph.Detail))                                          // 获取个人帖子详情
	postGroup.GET("/detail_pub/:postId", WrapParam(ph.DetailPub))                                   // 获取公开帖子详情
	postGroup.GET("/detail_post/:postId", casbinMiddleware.CheckCasbin(), WrapParam(ph.DetailPost)) // 管理员获取帖子详情

	// 统计相关
	postGroup.GET("/stats", casbinMiddleware.CheckCasbin(), WrapQuery(ph.GetPostCount)) // 获取帖子统计

	// 互动相关
	postGroup.POST("/like", WrapBody(ph.Like))       // 点赞/取消点赞
	postGroup.POST("/collect", WrapBody(ph.Collect)) // 收藏/取消收藏
}

// Edit 创建新帖子
func (ph *PostHandler) Edit(ctx *gin.Context, req req.EditReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)

	id, err := ph.svc.Create(ctx, domain.Post{
		ID:       req.PostId,
		Content:  req.Content,
		Title:    req.Title,
		PlateID:  req.PlateID,
		AuthorID: uc.Uid,
	})
	if err != nil {
		return Result{
			Code: PostEditERRORCode,
			Msg:  PostEditERROR,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostEditSuccess,
		Data: id,
	}, nil
}

// Update 更新帖子内容
func (ph *PostHandler) Update(ctx *gin.Context, req req.UpdateReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	
	if err := ph.svc.Update(ctx, domain.Post{
		ID:       req.PostId,
		Title:    req.Title,
		Content:  req.Content,
		PlateID:  req.PlateID,
		AuthorID: uc.Uid,
	}); err != nil {
		return Result{
			Code: PostUpdateERRORCode,
			Msg:  PostUpdateERROR,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostUpdateSuccess,
	}, nil
}

// Publish 发布帖子
func (ph *PostHandler) Publish(ctx *gin.Context, req req.PublishReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)

	if err := ph.svc.Publish(ctx, domain.Post{
		ID:       req.PostId,
		AuthorID: uc.Uid,
	}); err != nil {
		return Result{
			Code: PostPublishERRORCode,
			Msg:  PostPublishERROR,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostPublishSuccess,
		Data: req.PostId,
	}, nil
}

// Withdraw 撤回帖子
func (ph *PostHandler) Withdraw(ctx *gin.Context, req req.WithDrawReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)

	if err := ph.svc.Withdraw(ctx, domain.Post{
		ID:       req.PostId,
		AuthorID: uc.Uid,
	}); err != nil {
		return Result{
			Code: PostWithdrawERRORCode,
			Msg:  PostWithdrawERROR,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostWithdrawSuccess,
		Data: req.PostId,
	}, nil
}

// List 获取个人帖子列表
func (ph *PostHandler) List(ctx *gin.Context, req req.ListReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)

	du, err := ph.svc.ListPosts(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
		Uid:  uc.Uid,
	})
	if err != nil {
		return Result{
			Code: PostListERRORCode,
			Msg:  PostListERROR,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostListSuccess,
		Data: du,
	}, nil
}

// ListPub 获取公开帖子列表
func (ph *PostHandler) ListPub(ctx *gin.Context, req req.ListPubReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)

	du, err := ph.svc.ListPublishedPosts(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
		Uid:  uc.Uid,
	})
	if err != nil {
		return Result{
			Code: PostListPubERRORCode,
			Msg:  PostListPubERROR,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostListPubSuccess,
		Data: du,
	}, nil
}

// Detail 获取帖子详情
func (ph *PostHandler) Detail(ctx *gin.Context, req req.DetailReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)

	post, err := ph.svc.GetDraftsByAuthor(ctx, req.PostId, uc.Uid)
	if err != nil {
		return Result{
			Code: PostGetDetailERRORCode,
			Msg:  PostGetDetailERROR,
		}, nil
	}

	if post.Content == "" && post.Title == "" {
		return Result{
			Code: PostGetDetailERRORCode,
			Msg:  PostGetDetailERROR,
		}, nil
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostGetDetailSuccess,
		Data: post,
	}, nil
}

// DetailPub 获取公开帖子详情
func (ph *PostHandler) DetailPub(ctx *gin.Context, req req.DetailReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)

	post, err := ph.svc.GetPublishedPostById(ctx, req.PostId, uc.Uid)
	if err != nil {
		return Result{
			Code: PostGetPubDetailERRORCode,
			Msg:  PostGetPubDetailERROR,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostGetPubDetailSuccess,
		Data: post,
	}, nil
}

// DeletePost 删除帖子
func (ph *PostHandler) DeletePost(ctx *gin.Context, req req.DeleteReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)

	if err := ph.svc.Delete(ctx, req.PostId, uc.Uid); err != nil {
		return Result{
			Code: PostDeleteERRORCode,
			Msg:  PostDeleteERROR,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostDeleteSuccess,
		Data: req.PostId,
	}, nil
}

// Like 点赞/取消点赞
func (ph *PostHandler) Like(ctx *gin.Context, req req.LikeReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)

	var err error

	if req.Liked {
		err = ph.intSvc.Like(ctx, ph.biz, req.PostId, uc.Uid)
	} else {
		err = ph.intSvc.CancelLike(ctx, ph.biz, req.PostId, uc.Uid)
	}

	if err != nil {
		return Result{
			Code: PostLikedERRORCode,
			Msg:  PostLikedERROR,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostLikedSuccess,
		Data: req.PostId,
	}, nil
}

// Collect 收藏/取消收藏
func (ph *PostHandler) Collect(ctx *gin.Context, req req.CollectReq) (Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)

	var err error

	if req.Collectd {
		err = ph.intSvc.Collect(ctx, ph.biz, req.PostId, req.CollectId, uc.Uid)
	} else {
		err = ph.intSvc.CancelCollect(ctx, ph.biz, req.PostId, req.CollectId, uc.Uid)
	}

	if err != nil {
		return Result{
			Code: PostCollectERRORCode,
			Msg:  PostCollectERROR,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostCollectSuccess,
		Data: req.PostId,
	}, nil
}

// ListPost 管理员获取所有帖子列表
func (ph *PostHandler) ListPost(ctx *gin.Context, req req.ListPostReq) (Result, error) {
	du, err := ph.svc.ListAllPost(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return Result{
			Code: PostListERRORCode,
			Msg:  PostListERROR,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostListSuccess,
		Data: du,
	}, nil
}

// DetailPost 管理员获取帖子详情
func (ph *PostHandler) DetailPost(ctx *gin.Context, req req.DetailPostReq) (Result, error) {
	post, err := ph.svc.GetPost(ctx, req.PostId)
	if err != nil {
		return Result{
			Code: PostGetDetailERRORCode,
			Msg:  PostGetDetailERROR,
		}, err
	}

	return Result{
		Code: RequestsOK,
		Msg:  PostGetDetailSuccess,
		Data: post,
	}, nil
}

// GetPostCount 获取帖子统计数据
func (ph *PostHandler) GetPostCount(ctx *gin.Context, _ req.GetPostCountReq) (Result, error) {
	count, err := ph.svc.GetPostCount(ctx)
	if err != nil {
		return Result{
			Code: PostGetCountERRORCode,
			Msg:  PostGetCountERROR,
		}, err
	}
	
	return Result{
		Code: RequestsOK,
		Msg:  PostGetCountSuccess,
		Data: count,
	}, nil
}
