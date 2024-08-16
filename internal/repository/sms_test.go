package repository_test

import (
	"context"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/GoSimplicity/LinkMe/ioc"
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
	client := ioc.InitSms()
	repo := repository.NewSmsRepository(d, c, logger, client)
	if er := repo.SendCode(context.Background(), "xxx"); er != nil {
		fmt.Println(er)
		return
	}
}
