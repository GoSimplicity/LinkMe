package models

type Job struct {
	Id          int64  `gorm:"primaryKey,autoIncrement"`               // 任务ID，主键，自增
	Name        string `gorm:"type:varchar(128);unique"`               // 任务名称，唯一
	Executor    string `gorm:"type:varchar(255)"`                      // 执行者，执行该任务的实体
	Expression  string `gorm:"type:varchar(255)"`                      // 调度表达式，用于描述任务的调度时间
	Cfg         string `gorm:"type:text"`                              // 配置，任务的具体配置信息
	Status      int    `gorm:"type:int"`                               // 任务状态，用于标识任务当前的状态（如启用、禁用等）
	Version     int    `gorm:"type:int"`                               // 版本号，用于乐观锁控制并发更新
	NextTime    int64  `gorm:"index"`                                  // 下次执行时间，Unix时间戳
	CreateTime  int64  `gorm:"column:created_at;type:bigint;not null"` // 创建时间，Unix时间戳
	UpdatedTime int64  `gorm:"column:updated_at;type:bigint;not null"` // 更新时间，Unix时间戳
}
