package repository_test

import (
	"LinkMe/internal/repository"
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/service"
	"LinkMe/ioc"
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"testing"
)

func TestSendCode(t *testing.T) {
	configFile := pflag.String("config", "../../config/config.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	logger := ioc.InitLogger()
	d := dao.NewSmsDAO(ioc.InitDB(), logger)
	c := cache.NewSMSCache(ioc.InitRedis())
	repo := repository.NewSmsRepository(d, c)
	client := ioc.InitSms()
	smsService := service.NewSmsService(repo, logger, client, c)
	if er := smsService.SendCode(context.Background(), "xxx"); er != nil {
		fmt.Println(er)
		return
	}
}
