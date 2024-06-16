package main

import (
	"LinkMe/internal/domain/events"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type Cmd struct {
	server   *gin.Engine
	Cron     *cron.Cron
	consumer []events.Consumer
}
