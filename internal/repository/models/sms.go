package models

import (
	"github.com/ecodeclub/ekit/sqlx"
)

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

// VCodeSmsLog 表示用户认证操作的日志记录
type VCodeSmsLog struct {
	Id          int64  `gorm:"column:id;primaryKey;autoIncrement"`           // 自增ID
	SmsId       int64  `gorm:"column:sms_id"`                                // 短信类型ID
	SmsType     string `gorm:"column:sms_type"`                              // 短信类型
	Mobile      string `gorm:"column:mobile"`                                // 手机号
	VCode       string `gorm:"column:v_code"`                                // 验证码
	Driver      string `gorm:"column:driver"`                                // 服务商类型
	Status      int64  `gorm:"column:status"`                                // 发送状态，1为成功，0为失败
	StatusCode  string `gorm:"column:status_code"`                           // 状态码
	CreateTime  int64  `gorm:"column:created_at;type:bigint;not null"`       // 创建时间
	UpdatedTime int64  `gorm:"column:updated_at;type:bigint;not null;index"` // 更新时间
	DeletedTime int64  `gorm:"column:deleted_at;type:bigint;index"`          // 删除时间
}
