package models

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

type UserOperationLog struct {
	Phone      string `gorm:"column:phone"`
	Action     string `gorm:"column:action"`
	CreateTime int64  `gorm:"column:create_Time"`
}
