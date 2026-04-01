//go:build integration
// +build integration

package repository_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/GoSimplicity/LinkMe/ioc"
	"github.com/spf13/viper"
)

func TestSendCode(t *testing.T) {
	viper.Reset()
	viper.Set("log.dir", filepath.Join(os.TempDir(), "linkme-test-logs"))
	viper.Set("db.dsn", os.Getenv("LINKME_DB_DSN"))
	viper.Set("redis.addr", os.Getenv("LINKME_REDIS_ADDR"))
	viper.Set("redis.password", os.Getenv("LINKME_REDIS_PASSWORD"))
	viper.Set("sms.provider", os.Getenv("LINKME_SMS_PROVIDER"))
	viper.Set("sms.tencent.smsID", os.Getenv("LINKME_SMS_TENCENT_SMSID"))
	viper.Set("sms.tencent.sign", os.Getenv("LINKME_SMS_TENCENT_SIGN"))
	viper.Set("sms.tencent.templateID", os.Getenv("LINKME_SMS_TENCENT_TEMPLATEID"))
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
