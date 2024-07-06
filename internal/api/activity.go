package api

import (
	"LinkMe/internal/constants"
	"LinkMe/internal/service"
	"LinkMe/middleware"
	. "LinkMe/pkg/ginp"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ActivityHandler struct {
	ce  *casbin.Enforcer
	svc service.ActivityService
	l   *zap.Logger
}

func NewActivityHandler(svc service.ActivityService, ce *casbin.Enforcer, l *zap.Logger) *ActivityHandler {
	return &ActivityHandler{
		svc: svc,
		ce:  ce,
		l:   l,
	}
}

func (ah *ActivityHandler) RegisterRoutes(server *gin.Engine) {
	casbinMiddleware := middleware.NewCasbinMiddleware(ah.ce, ah.l)
	historyGroup := server.Group("/api/activity")
	historyGroup.GET("/recent", casbinMiddleware.CheckCasbin(), WrapQuery(ah.GetRecentActivity)) // 获取最近的活动记录
}

func (ah *ActivityHandler) GetRecentActivity(ctx *gin.Context, _ GetRecentActivityReq) (Result, error) {
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
