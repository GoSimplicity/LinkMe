package domain

type Plate struct {
	ID          int64  // 板块ID
	Name        string // 板块名称
	Description string // 板块描述
	Uid         int64  // 操作人
	CreatedAt   int64  // 创建时间
	UpdatedAt   int64  // 更新时间
	DeletedAt   int64  // 删除时间
	Deleted     bool   // 删除状态
}
