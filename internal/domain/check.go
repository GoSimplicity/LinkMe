package domain

type Check struct {
	ID        int64  // 审核ID
	PostID    int64  // 帖子ID
	Content   string // 审核内容
	Title     string // 审核标签
	UserID    int64  // 提交审核的用户ID
	Status    string // 审核状态
	Remark    string // 审核备注
	CreatedAt int64  // 创建时间
	UpdatedAt int64  // 更新时间
}

type CheckList struct {
	ID        int64  // 审核ID
	PostID    int64  // 帖子ID
	Title     string // 审核标签
	UserID    int64  // 提交审核的用户ID
	Status    string // 审核状态
	Remark    string // 审核备注
	CreatedAt int64  // 创建时间
	UpdatedAt int64  // 更新时间
}
