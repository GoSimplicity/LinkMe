package models

import "github.com/ecodeclub/ekit/sqlx"

type Sms struct {
	Id       int64                      // 唯一标识
	Config   sqlx.JsonColumn[SmsConfig] // SMS 配置，存储为 JSON 格式
	RetryCnt int                        // 当前重试次数
	RetryMax int                        // 最大重试次数
	Status   uint8                      // 状态码
	Ctime    int64                      // 创建时间
	Utime    int64                      `gorm:"index"` // 更新时间，添加索引
}

// SmsConfig 表示 SMS 配置的结构
type SmsConfig struct {
	TplId   string   // 模板ID
	Args    []string // 模板参数
	Numbers []string // 接收短信的手机号列表
}
