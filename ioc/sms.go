package ioc

import (
	"LinkMe/pkg/sms"
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	tencentsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func newClient() *tencentsms.Client {
	secretId := viper.GetString("sms.tencent.secretId")
	secretKey := viper.GetString("sms.tencent.secretKey")
	endPoint := viper.GetString("sms.tencent.endPoint")
	credential := common.NewCredential(
		secretId,
		secretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = endPoint
	client, _ := tencentsms.NewClient(credential, regions.Nanjing, cpf)
	return client
}

func InitSms() *sms.TencentSms {
	smsID := viper.GetString("sms.tencent.smsID")
	sign := viper.GetString("sms.tencent.sign")
	templateID := viper.GetString("sms.tencent.templateID")
	return sms.NewTencentSms(newClient(), smsID, sign, templateID)
}
