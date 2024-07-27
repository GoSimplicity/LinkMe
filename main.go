package main

import (
	"context"
	"github.com/GoSimplicity/LinkMe/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	Init()
}

func Init() {
	config.InitViper()
	cmd := InitWebServer()
	server := cmd.server
	server.GET("/headers", func(c *gin.Context) {
		headers := c.Request.Header
		// 打印所有请求头
		for key, values := range headers {
			for _, value := range values {
				c.String(http.StatusOK, "%s: %s\n", key, value)
			}
		}
	})
	for _, s := range cmd.consumer {
		err := s.Start(context.Background())
		if err != nil {
			panic(err)
		}
	}
	cmd.Cron.Start() // 启动定时任务
	if er := server.Run(":9999"); er != nil {
		panic(er)
	}
}
