package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AccessLog struct {
	Path     string        `json:"path"`     // 请求路径
	Method   string        `json:"method"`   // 请求方法
	ReqBody  string        `json:"reqBody"`  // 请求体内容
	Status   int           `json:"status"`   // 响应状态码
	RespBody string        `json:"respBody"` // 响应体内容
	Duration time.Duration `json:"duration"` // 请求处理耗时
}

type LogMiddleware struct {
	l *zap.Logger
}

func NewLogMiddleware(l *zap.Logger) *LogMiddleware {
	return &LogMiddleware{
		l: l,
	}
}

// Log 日志中间件
func (lm *LogMiddleware) Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			lm.l.Error("请求体读取失败", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// 由于读取请求体会消耗掉c.Request.Body，所以需要重新设置回上下文
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		al := AccessLog{
			Path:    path,
			Method:  method,
			ReqBody: string(bodyBytes),
		}
		c.Next()
		// 记录响应状态码和响应体
		al.Status = c.Writer.Status()
		al.RespBody = c.Writer.Header().Get("Content-Type")
		al.Duration = time.Since(start)
		lm.l.Info("请求日志", zap.Any("accessLog", al))
	}
}
