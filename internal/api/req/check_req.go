package req

// SubmitCheckReq 定义了提交审核请求的结构体
type SubmitCheckReq struct {
	PostID  int64  `json:"postId" binding:"required"`  // 帖子ID
	Content string `json:"content" binding:"required"` // 审核内容
	Title   string `json:"title" binding:"required"`   // 审核标题
	UserID  int64  `json:"userId" binding:"required"`  // 提交审核的用户ID
}

// ApproveCheckReq 定义了审核通过请求的结构体
type ApproveCheckReq struct {
	CheckID int64  `json:"checkId" binding:"required"` // 审核ID
	Remark  string `json:"remark"`                     // 审核通过备注
}

// RejectCheckReq 定义了审核拒绝请求的结构体
type RejectCheckReq struct {
	CheckID int64  `json:"checkId" binding:"required"` // 审核ID
	UserID  int64  `json:"userId" binding:"required"`  // 审核拒绝的用户ID
	Remark  string `json:"remark" binding:"required"`  // 审核拒绝原因
}

// ListCheckReq 定义了获取审核列表请求的结构体
type ListCheckReq struct {
	Page int    `form:"page" binding:"required"` // 页码
	Size *int64 `form:"size" binding:"required"` // 每页数量
}

// CheckDetailReq 定义了获取审核详情请求的结构体
type CheckDetailReq struct {
	CheckID int64 `json:"checkId" binding:"required"` // 审核ID
}

type GetCheckCount struct {
}
