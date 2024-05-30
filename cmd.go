package main

import (
	"LinkMe/internal/domain/events"
	"github.com/gin-gonic/gin"
)

type Cmd struct {
	server   *gin.Engine
	consumer []events.Consumer
}
