package ginp

import (
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
	"github.com/gin-gonic/gin"
	"net/http"
)

// WrapBody 是一个中间件，用于包裹业务逻辑函数，自动绑定请求体、处理响应并集中管理错误处理。
func WrapBody[Req any](bizFn func(ctx *gin.Context, req Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.ShouldBindJSON(&req); err != nil {
			apiresponse.BadRequestError(ctx, "无效的请求参数")
			return
		}

		res, err := bizFn(ctx, req)
		if err != nil {
			apiresponse.InternalServerErrorWithDetails(ctx, gin.H{"error": err.Error()}, res.Msg)
			return
		}

		ctx.JSON(http.StatusOK, apiresponse.ApiResponse{
			Code:    res.Code,
			Data:    res.Data,
			Message: res.Msg,
		})
	}
}

func WrapParam[Req any](bizFn func(ctx *gin.Context, req Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.ShouldBindUri(&req); err != nil {
			apiresponse.BadRequestError(ctx, "无效的路径参数")
			return
		}

		res, err := bizFn(ctx, req)
		if err != nil {
			apiresponse.InternalServerErrorWithDetails(ctx, gin.H{"error": err.Error()}, res.Msg)
			return
		}
		ctx.JSON(http.StatusOK, apiresponse.ApiResponse{
			Code:    res.Code,
			Data:    res.Data,
			Message: res.Msg,
		})
	}
}

func WrapQuery[Req any](bizFn func(ctx *gin.Context, req Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.ShouldBindQuery(&req); err != nil {
			apiresponse.BadRequestError(ctx, "无效的查询参数")
			return
		}
		res, err := bizFn(ctx, req)
		if err != nil {
			apiresponse.InternalServerErrorWithDetails(ctx, gin.H{"error": err.Error()}, res.Msg)
			return
		}
		ctx.JSON(http.StatusOK, apiresponse.ApiResponse{
			Code:    res.Code,
			Data:    res.Data,
			Message: res.Msg,
		})
	}
}

// WrapNoParam 是一个中间件，用于包裹不需要请求参数的业务逻辑函数，处理响应并集中管理错误处理。
func WrapNoParam(bizFn func(ctx *gin.Context) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := bizFn(ctx)
		if err != nil {
			apiresponse.InternalServerErrorWithDetails(ctx, gin.H{"error": err.Error()}, res.Msg)
			return
		}
		ctx.JSON(http.StatusOK, apiresponse.ApiResponse{
			Code:    res.Code,
			Data:    res.Data,
			Message: res.Msg,
		})
	}
}
