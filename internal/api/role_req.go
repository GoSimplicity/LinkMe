package api

type CreateRoleReq struct {
	Name string `json:"name"`
}

type CreatePermissionReq struct {
	Name string `json:"name"`
}

type AssignPermissionToRoleReq struct {
	RoleID       int64 `json:"roleId"`
	PermissionID int64 `json:"permissionId"`
}
