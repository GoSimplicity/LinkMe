package qqEmail

import (
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) (err error) {
	from := viper.GetString("email.qq.from")
	password := viper.GetString("email.qq.password")
	host := viper.GetString("email.qq.host")
	port := viper.GetInt("email.qq.port")
	if host == "" {
		host = "smtp.qq.com"
	}
	if port == 0 {
		port = 587
	}

	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// 使用QQ邮箱SMTP服务器
	d := gomail.NewDialer(host, port, from, password)

	// 发送邮件
	return d.DialAndSend(m)
}
