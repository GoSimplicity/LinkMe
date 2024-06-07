package domain

const (
	Pending  = "Pending"  // 待审核状态
	Approved = "Approved" // 审核通过状态
	Rejected = "Rejected" // 审核拒绝状态
)

type Check struct {
	ID        int64  // 审核ID
	PostID    int64  // 帖子ID
	Content   string // 审核内容
	Title     string // 审核标签
	UserID    int64  // 提交审核的用户ID
	Status    string // 审核状态
	Remark    string // 审核备注
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}
