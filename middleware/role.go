package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware 创建一个 Casbin 中间件
func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户身份、请求的 URL 和请求方法
		sub := c.GetString("username")
		obj := c.Request.URL.Path
		act := c.Request.Method
		// 使用 Casbin 检查权限
		ok, err := enforcer.Enforce(sub, obj, act)
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
