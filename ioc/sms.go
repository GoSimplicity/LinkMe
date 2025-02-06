package ioc

import (
	"github.com/GoSimplicity/LinkMe/pkg/sms"
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	tencentsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// newClient 创建腾讯云短信客户端
func newClient() *tencentsms.Client {
	secretId := viper.GetString("sms.tencent.secretId")
	secretKey := viper.GetString("sms.tencent.secretKey")
	endPoint := viper.GetString("sms.tencent.endPoint")
	credential := common.NewCredential(
		secretId,
		secretKey,
	)
	// 创建客户端配置
	cpf := profile.NewClientProfile()
	// 设置腾讯云短信服务端点
	cpf.HttpProfile.Endpoint = endPoint
	// 创建腾讯云短信客户端
	client, _ := tencentsms.NewClient(credential, regions.Nanjing, cpf)
	return client
}

// InitSms 初始化腾讯云短信实例
func InitSms() *sms.TencentSms {
	smsID := viper.GetString("sms.tencent.smsID")
	sign := viper.GetString("sms.tencent.sign")
	templateID := viper.GetString("sms.tencent.templateID")
	// 创建腾讯云短信实例
	return sms.NewTencentSms(newClient(), smsID, sign, templateID)
}
