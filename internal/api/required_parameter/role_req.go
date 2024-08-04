package required_parameter

// ListPermissionsReq 是获取权限列表的请求
type ListPermissionsReq struct {
}

// AssignPermissionReq 是分配权限的请求
type AssignPermissionReq struct {
	UserName string `json:"userName,omitempty"` // 用户
	Path     string `json:"path,omitempty"`     // 接口路径
	Method   string `json:"method,omitempty"`   // 操作类型
}

// AssignPermissionRoleReq 是分配角色的请求
type AssignPermissionRoleReq struct {
	UserName string `json:"userName,omitempty"` // 用户
	RoleName string `json:"roleName,omitempty"` // 角色
}

// RemovePermissionReq 是移除权限的请求
type RemovePermissionReq struct {
	UserName string `json:"userName,omitempty"` // 用户
	Path     string `json:"path,omitempty"`     // 接口路径
	Method   string `json:"method,omitempty"`   // 操作类型
}

// RemovePermissionRoleReq 是移除角色的请求
type RemovePermissionRoleReq struct {
	UserName string `json:"userName,omitempty"` // 用户
	RoleName string `json:"roleName,omitempty"` // 角色
}
