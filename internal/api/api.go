package api

import (
	"strconv"

	"github.com/GoSimplicity/LinkMe/internal/api/req"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApiHandler struct {
	svc service.ApiService
	l   *zap.Logger
}

func NewApiHandler(svc service.ApiService, l *zap.Logger) *ApiHandler {
	return &ApiHandler{
		svc: svc,
		l:   l,
	}
}

func (h *ApiHandler) RegisterRoutes(server *gin.Engine) {
	apiGroup := server.Group("/api/apis")

	apiGroup.POST("/list", h.ListApis)
	apiGroup.POST("/create", h.CreateAPI)
	apiGroup.POST("/update", h.UpdateAPI)
	apiGroup.DELETE("/:id", h.DeleteAPI)
}

// ListApis 获取API列表
func (a *ApiHandler) ListApis(c *gin.Context) {
	var req req.ListApisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 调用service层获取API列表
	apis, total, err := a.svc.ListApis(c.Request.Context(), req.PageNumber, req.PageSize)
	if err != nil {
		a.l.Error("获取API列表失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.SuccessWithData(c, gin.H{
		"list":  apis,
		"total": total,
	})
}

// CreateAPI 创建新的API
func (a *ApiHandler) CreateAPI(c *gin.Context) {
	var req req.CreateApiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 构建API对象
	api := &domain.Api{
		Name:        req.Name,
		Path:        req.Path,
		Method:      req.Method,
		Description: req.Description,
		Version:     req.Version,
		Category:    req.Category,
		IsPublic:    req.IsPublic,
	}

	if err := a.svc.CreateApi(c.Request.Context(), api); err != nil {
		a.l.Error("创建API失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

// UpdateAPI 更新API信息
func (a *ApiHandler) UpdateAPI(c *gin.Context) {
	var r req.UpdateApiRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		a.l.Error("绑定请求参数失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	// 构建更新的API对象
	api := &domain.Api{
		ID:          r.Id,
		Name:        r.Name,
		Path:        r.Path,
		Method:      r.Method,
		Description: r.Description,
		Version:     r.Version,
		Category:    r.Category,
		IsPublic:    r.IsPublic,
	}

	if err := a.svc.UpdateApi(c.Request.Context(), api); err != nil {
		a.l.Error("更新API失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}

// DeleteAPI 删除API
func (a *ApiHandler) DeleteAPI(c *gin.Context) {
	// 从URL参数中获取API ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		a.l.Error("解析ID失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	if err := a.svc.DeleteApi(c.Request.Context(), id); err != nil {
		a.l.Error("删除API失败", zap.Error(err))
		apiresponse.Error(c)
		return
	}

	apiresponse.Success(c)
}
