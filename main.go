package main

import (
	"LinkMe/config"
	"context"
)

func main() {
	Init()
}

func Init() {
	config.InitViper()
	cmd := InitWebServer()
	server := cmd.server
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
