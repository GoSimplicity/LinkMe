package sms

import (
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencent "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// TencentSms 腾讯云SMS
type TencentSms struct {
	client     *tencent.Client
	smsID      string
	signName   string
	TemplateID string
}

func NewTencentSms(client *tencent.Client, smsID string, signName string, TemplateID string) *TencentSms {
	return &TencentSms{
		client:     client,
		smsID:      smsID,
		signName:   signName,
		TemplateID: TemplateID,
	}
}

// Send 发送短信
func (s *TencentSms) Send(ctx context.Context, args []string, numbers ...string) (smsID string, driver string, err error) {
	request := tencent.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = common.StringPtr(s.smsID)
	request.SignName = common.StringPtr(s.signName)
	request.TemplateId = common.StringPtr(s.TemplateID)

	request.TemplateParamSet = common.StringPtrs(args)
	request.PhoneNumberSet = common.StringPtrs(numbers)

	response, err := s.client.SendSms(request)
	if err != nil {
		return s.smsID, "tencent", err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			continue
		}
		status := *statusPtr
		if status.Code == nil || *(status.Code) != "Ok" {
			// 发送失败
			return s.smsID, "tencent", fmt.Errorf("send sms messages failed，code: %s, msg: %s", *status.Code, *status.Message)
		}
	}
	return s.smsID, "tencent", nil
}
