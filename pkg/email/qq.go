package qqEmail

import (
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) (err error) {
	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", "xxx@qq.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// 使用QQ邮箱SMTP服务器
	d := gomail.NewDialer("smtp.qq.com", 587, "xxx@qq.com", "xxx")

	// 发送邮件
	return d.DialAndSend(m)
}
