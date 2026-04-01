package api

import (
	"github.com/GoSimplicity/LinkMe/pkg/apiresponse"
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"github.com/gin-gonic/gin"
)

func requireUser(ctx *gin.Context) (ijwt.UserClaims, bool) {
	user, exists := ctx.Get("user")
	if !exists {
		apiresponse.UnauthorizedErrorWithDetails(ctx, nil, "未登录或登录已过期")
		return ijwt.UserClaims{}, false
	}

	claims, ok := user.(ijwt.UserClaims)
	if !ok || claims.Uid == 0 {
		apiresponse.UnauthorizedErrorWithDetails(ctx, nil, "未登录或登录已过期")
		return ijwt.UserClaims{}, false
	}

	return claims, true
}

func currentUserID(ctx *gin.Context) int64 {
	user, exists := ctx.Get("user")
	if !exists {
		return 0
	}

	claims, ok := user.(ijwt.UserClaims)
	if !ok {
		return 0
	}

	return claims.Uid
}
