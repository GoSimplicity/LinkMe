package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/req"
	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	svc    service.PostService
	intSvc service.InteractiveService
}

func NewPostHandler(svc service.PostService, intSvc service.InteractiveService) *PostHandler {
	return &PostHandler{
		svc:    svc,
		intSvc: intSvc,
	}
}

func (ph *PostHandler) RegisterRoutes(server *gin.Engine) {
	postGroup := server.Group("/api/posts")

	postGroup.POST("/edit", ph.Edit)                   // 编辑帖子
	postGroup.POST("/update", ph.Update)               // 更新帖子
	postGroup.POST("/publish", ph.Publish)             // 发布帖子
	postGroup.POST("/withdraw", ph.Withdraw)           // 撤回帖子
	postGroup.DELETE("/delete/:postId", ph.DeletePost) // 删除帖子
	postGroup.POST("/list", ph.List)                   // 获取个人帖子列表
	postGroup.POST("/list_pub", ph.ListPub)            // 获取公开帖子列表
	postGroup.POST("/list_all", ph.ListAll)            // 获取所有帖子列表
	postGroup.GET("/get/:id", ph.GetPost)              // 获取帖子详情
	postGroup.GET("/detail/:postId", ph.Detail)        // 获取个人帖子详情
	postGroup.GET("/detail_pub/:postId", ph.DetailPub) // 获取公开帖子详情
	postGroup.POST("/like", ph.Like)                   // 点赞/取消点赞
	postGroup.POST("/collect", ph.Collect)             // 收藏/取消收藏
}

// Edit 创建新帖子
func (ph *PostHandler) Edit(ctx *gin.Context) {
	var req req.EditReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	id, err := ph.svc.Create(ctx, domain.Post{
		ID:      req.PostId,
		Content: req.Content,
		Title:   req.Title,
		PlateID: req.PlateID,
		Uid:     uc.Uid,
	})
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	apiresponse.SuccessWithData(ctx, id)
}

// Update 更新帖子内容
func (ph *PostHandler) Update(ctx *gin.Context) {
	var req req.UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	if err := ph.svc.Update(ctx, domain.Post{
		ID:      req.PostId,
		Content: req.Content,
		Title:   req.Title,
		PlateID: req.PlateID,
		Uid:     uc.Uid,
	}); err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	apiresponse.Success(ctx)
}

// Publish 发布帖子
func (ph *PostHandler) Publish(ctx *gin.Context) {
	var req req.PublishReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	if err := ph.svc.Publish(ctx, req.PostId, uc.Uid); err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	apiresponse.SuccessWithData(ctx, req.PostId)
}

// Withdraw 撤回帖子
func (ph *PostHandler) Withdraw(ctx *gin.Context) {
	var req req.WithDrawReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	if err := ph.svc.Withdraw(ctx, req.PostId, uc.Uid); err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	apiresponse.SuccessWithData(ctx, req.PostId)
}

// List 获取个人帖子列表
func (ph *PostHandler) List(ctx *gin.Context) {
	var req req.ListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	du, err := ph.svc.ListPosts(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
		Uid:  uc.Uid,
	})
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	apiresponse.SuccessWithData(ctx, du)
}

// ListPub 获取公开帖子列表
func (ph *PostHandler) ListPub(ctx *gin.Context) {
	var req req.ListPubReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	du, err := ph.svc.ListPublishPosts(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
		Uid:  uc.Uid,
	})
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	apiresponse.SuccessWithData(ctx, du)
}

// Detail 获取帖子详情
func (ph *PostHandler) Detail(ctx *gin.Context) {
	var req req.DetailReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	post, err := ph.svc.GetPostById(ctx, req.PostId, uc.Uid)
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	if post.Content == "" && post.Title == "" {
		apiresponse.ErrorWithMessage(ctx, PostGetDetailERROR)
		return
	}

	apiresponse.SuccessWithData(ctx, post)
}

// DetailPub 获取公开帖子详情
func (ph *PostHandler) DetailPub(ctx *gin.Context) {
	var req req.DetailReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	post, err := ph.svc.GetPublishPostById(ctx, req.PostId, uc.Uid)
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	apiresponse.SuccessWithData(ctx, post)
}

// DeletePost 删除帖子
func (ph *PostHandler) DeletePost(ctx *gin.Context) {
	var req req.DeleteReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	if err := ph.svc.Delete(ctx, req.PostId, uc.Uid); err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	apiresponse.SuccessWithData(ctx, req.PostId)
}

// Like 点赞/取消点赞
func (ph *PostHandler) Like(ctx *gin.Context) {
	var req req.LikeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	var err error

	if req.Liked {
		err = ph.intSvc.Like(ctx, req.PostId, uc.Uid)
	} else {
		err = ph.intSvc.CancelLike(ctx, req.PostId, uc.Uid)
	}

	if err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	apiresponse.SuccessWithData(ctx, req.PostId)
}

// Collect 收藏/取消收藏
func (ph *PostHandler) Collect(ctx *gin.Context) {
	var req req.CollectReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)

	var err error

	if req.Collectd {
		err = ph.intSvc.Collect(ctx, req.PostId, uc.Uid)
	} else {
		err = ph.intSvc.CancelCollect(ctx, req.PostId, uc.Uid)
	}

	if err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	apiresponse.SuccessWithData(ctx, req.PostId)
}

// ListAll 获取所有帖子列表
func (ph *PostHandler) ListAll(ctx *gin.Context) {
	var req req.ListPostReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	posts, err := ph.svc.ListAll(ctx, domain.Pagination{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	apiresponse.SuccessWithData(ctx, posts)
}

// GetPost 获取帖子详情
func (ph *PostHandler) GetPost(ctx *gin.Context) {
	var req req.DetailPostReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		apiresponse.ErrorWithMessage(ctx, "无效的请求参数")
		return
	}

	post, err := ph.svc.GetPost(ctx, req.PostId)
	if err != nil {
		apiresponse.ErrorWithMessage(ctx, err.Error())
		return
	}

	apiresponse.SuccessWithData(ctx, post)
}
