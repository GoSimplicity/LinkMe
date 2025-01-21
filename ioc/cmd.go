package ioc

import (
	"github.com/GoSimplicity/LinkMe/internal/domain/events"
	"github.com/GoSimplicity/LinkMe/internal/job"
	"github.com/GoSimplicity/LinkMe/internal/mock"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type Cmd struct {
	Server    *gin.Engine
	Consumer  []events.Consumer
	Mock      mock.MockUserRepository
	Routes    *job.Routes
	Asynq     *asynq.Server
	Scheduler *job.TimedScheduler
}
