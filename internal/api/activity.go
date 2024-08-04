package api

import (
	"github.com/GoSimplicity/LinkMe/internal/api/required_parameter"
	"github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/middleware"
	. "github.com/GoSimplicity/LinkMe/pkg/ginp"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	ce  *casbin.Enforcer
	svc service.ActivityService
}

func NewActivityHandler(svc service.ActivityService, ce *casbin.Enforcer) *ActivityHandler {
	return &ActivityHandler{
		svc: svc,
		ce:  ce,
	}
}

func (ah *ActivityHandler) RegisterRoutes(server *gin.Engine) {
	casbinMiddleware := middleware.NewCasbinMiddleware(ah.ce)
	historyGroup := server.Group("/api/activity")
	historyGroup.GET("/recent", casbinMiddleware.CheckCasbin(), WrapQuery(ah.GetRecentActivity)) // 获取最近的活动记录
}

// GetRecentActivity 获取最近的活动记录
func (ah *ActivityHandler) GetRecentActivity(ctx *gin.Context, _ required_parameter.GetRecentActivityReq) (Result, error) {
	activity, err := ah.svc.GetRecentActivity(ctx)
	if err != nil {
		return Result{
			Code: constants.GetRecentActivityERRORCode,
			Msg:  constants.GetRecentActivityERROR,
		}, err
	}
	return Result{
		Code: constants.RequestsOK,
		Msg:  constants.GetCheckSuccess,
		Data: activity,
	}, nil
}
