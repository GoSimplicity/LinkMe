package middleware

import (
	ijwt "github.com/GoSimplicity/LinkMe/utils/jwt"
	"net/http"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type CasbinMiddleware struct {
	enforcer *casbin.Enforcer
}

func NewCasbinMiddleware(enforcer *casbin.Enforcer) *CasbinMiddleware {
	return &CasbinMiddleware{
		enforcer: enforcer,
	}
}

// CheckCasbin 创建一个 Casbin 中间件
func (cm *CasbinMiddleware) CheckCasbin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户身份
		userClaims, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
			c.Abort()
			return
		}
		sub, ok := userClaims.(ijwt.UserClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid user claims"})
			c.Abort()
			return
		}
		if sub.Uid == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid user ID"})
			c.Abort()
			return
		}
		// 将用户ID转换为字符串
		userIDStr := strconv.FormatInt(sub.Uid, 10)
		// 获取请求的 URL 和请求方法
		obj := c.Request.URL.Path
		act := c.Request.Method
		// 使用 Casbin 检查权限
		ok, err := cm.enforcer.Enforce(userIDStr, obj, act)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error occurred when enforcing policy"})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"message": "You don't have permission to access this resource"})
			c.Abort()
			return
		}
		// 继续处理请求
		c.Next()
	}
}
