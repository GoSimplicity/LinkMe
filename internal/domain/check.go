package domain

const (
	UnderReview uint8 = iota
	Approved
	UnApproved
)

type Check struct {
	ID        int64  // 审核ID
	PostID    uint   // 帖子ID
	Content   string // 审核内容
	Title     string // 审核标签
	Uid       int64  // 提交审核的用户ID
	PlateID   int64  // 板块id
	Status    uint8  // 审核状态
	Remark    string // 审核备注
	CreatedAt int64  // 创建时间
	UpdatedAt int64  // 更新时间
}
