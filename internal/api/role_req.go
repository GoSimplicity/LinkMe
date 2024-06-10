package api

// ListPermissionsReq 是获取权限列表的请求
type ListPermissionsReq struct {
	UserID int64 `json:"userId,omitempty"` // 用户ID
}

// AssignPermissionReq 是分配权限的请求
type AssignPermissionReq struct {
	UserID int64  `json:"userId,omitempty"` // 用户ID
	Path   string `json:"path,omitempty"`   // 接口路径
	Method string `json:"method,omitempty"` // 操作类型
}

// RemovePermissionReq 是移除权限的请求
type RemovePermissionReq struct {
	UserID int64  `json:"userId,omitempty"` // 用户ID
	Path   string `json:"path,omitempty"`   // 接口路径
	Method string `json:"method,omitempty"` // 操作类型
}
