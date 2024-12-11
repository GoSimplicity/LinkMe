package req

// API相关
type CreateApiRequest struct {
	Name        string `json:"name" binding:"required"`       // API名称
	Path        string `json:"path" binding:"required"`       // API路径
	Method      int    `json:"method" binding:"required"`     // 请求方法
	Description string `json:"description"`                   // API描述
	Version     string `json:"version"`                       // API版本
	Category    int    `json:"category"`                      // API分类
	IsPublic    int    `json:"is_public" binding:"oneof=0 1"` // 是否公开
}

type GetApiRequest struct {
	Id int `json:"id" binding:"required,gt=0"` // API ID
}

type UpdateApiRequest struct {
	Id          int64  `json:"id" binding:"required,gt=0"`    // API ID
	Name        string `json:"name" binding:"required"`       // API名称
	Path        string `json:"path" binding:"required"`       // API路径
	Method      int    `json:"method" binding:"required"`     // 请求方法
	Description string `json:"description"`                   // API描述
	Version     string `json:"version"`                       // API版本
	Category    int    `json:"category"`                      // API分类
	IsPublic    int    `json:"is_public" binding:"oneof=0 1"` // 是否公开
}

type DeleteApiRequest struct {
	Id int `json:"id" binding:"required,gt=0"` // API ID
}

type ListApisRequest struct {
	PageNumber int `json:"page_number" binding:"required,gt=0"` // 页码
	PageSize   int `json:"page_size" binding:"required,gt=0"`   // 每页数量
}

// 菜单相关
type CreateMenuRequest struct {
	Name      string `json:"name" binding:"required"`    // 菜单名称
	Path      string `json:"path" binding:"required"`    // 菜单路径
	ParentId  int    `json:"parent_id" binding:"gte=0"`  // 父菜单ID
	Component string `json:"component"`                  // 组件
	Icon      string `json:"icon"`                       // 图标
	SortOrder int    `json:"sort_order" binding:"gte=0"` // 排序
	RouteName string `json:"route_name"`                 // 路由名称
	Hidden    int    `json:"hidden" binding:"oneof=0 1"` // 是否隐藏
}

type GetMenuRequest struct {
	Id int `json:"id" binding:"required,gt=0"` // 菜单ID
}

type UpdateMenuRequest struct {
	Id        int    `json:"id" binding:"required,gt=0"` // 菜单ID
	Name      string `json:"name" binding:"required"`    // 菜单名称
	Path      string `json:"path" binding:"required"`    // 菜单路径
	ParentId  int    `json:"parent_id" binding:"gte=0"`  // 父菜单ID
	Component string `json:"component"`                  // 组件
	Icon      string `json:"icon"`                       // 图标
	SortOrder int    `json:"sort_order" binding:"gte=0"` // 排序
	RouteName string `json:"route_name"`                 // 路由名称
	Hidden    int    `json:"hidden" binding:"oneof=0 1"` // 是否隐藏
}

type DeleteMenuRequest struct {
	Id int `json:"id" binding:"required,gt=0"` // 菜单ID
}

type ListMenusRequest struct {
	PageNumber int  `json:"page_number" binding:"required,gt=0"` // 页码
	PageSize   int  `json:"page_size" binding:"required,gt=0"`   // 每页数量
	IsTree     bool `json:"is_tree"`                             // 是否树形结构
}

// 角色相关
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`        // 角色名称
	Description string `json:"description"`                    // 角色描述
	RoleType    int    `json:"role_type" binding:"required"`   // 角色类型
	IsDefault   int    `json:"is_default" binding:"oneof=0 1"` // 是否默认角色
}

type GetRoleRequest struct {
	Id int `json:"id" binding:"required,gt=0"` // 角色ID
}

type UpdateRoleRequest struct {
	Id          int64  `json:"id" binding:"required,gt=0"`     // 角色ID
	Name        string `json:"name" binding:"required"`        // 角色名称
	Description string `json:"description"`                    // 角色描述
	RoleType    int    `json:"role_type" binding:"required"`   // 角色类型
	IsDefault   int    `json:"is_default" binding:"oneof=0 1"` // 是否默认角色
}

type DeleteRoleRequest struct {
	Id int `json:"id" binding:"required,gt=0"` // 角色ID
}

type ListRolesRequest struct {
	PageNumber int `json:"page_number" binding:"required,gt=0"` // 页码
	PageSize   int `json:"page_size" binding:"required,gt=0"`   // 每页数量
}

type AssignPermissionsRequest struct {
	RoleId  int   `json:"role_id" binding:"required,gt=0"` // 角色ID
	MenuIds []int `json:"menu_ids"`                        // 菜单ID列表
	ApiIds  []int `json:"api_ids"`                         // API ID列表
}

type AssignRoleToUserRequest struct {
	UserId  int   `json:"user_id" binding:"required,gt=0"` // 用户ID
	RoleIds []int `json:"role_ids" binding:"required"`     // 角色ID列表
}

type RemoveRoleFromUserRequest struct {
	UserId  int   `json:"user_id" binding:"required,gt=0"` // 用户ID
	RoleIds []int `json:"role_ids" binding:"required"`     // 角色ID列表
}

type RemoveUserPermissionsRequest struct {
	UserId int `json:"user_id" binding:"required,gt=0"` // 用户ID
}

type RemoveRoleApiPermissionsRequest struct {
	RoleIds []int `json:"role_ids" binding:"required"` // 角色ID列表
	ApiIds  []int `json:"api_ids" binding:"required"`  // API ID列表
}

type RemoveRoleMenuPermissionsRequest struct {
	RoleIds []int `json:"role_ids" binding:"required"` // 角色ID列表
	MenuIds []int `json:"menu_ids" binding:"required"` // 菜单ID列表
}

type RemoveUserApiPermissionsRequest struct {
	UserId int   `json:"user_id" binding:"required,gt=0"` // 用户ID
	ApiIds []int `json:"api_ids" binding:"required"`      // API ID列表
}

type RemoveUserMenuPermissionsRequest struct {
	UserId  int   `json:"user_id" binding:"required,gt=0"` // 用户ID
	MenuIds []int `json:"menu_ids" binding:"required"`     // 菜单ID列表
}

type AssignApiPermissionsToUserRequest struct {
	UserId int   `json:"user_id" binding:"required,gt=0"` // 用户ID
	ApiIds []int `json:"api_ids" binding:"required"`      // API ID列表
}

type AssignMenuPermissionsToUserRequest struct {
	UserId  int   `json:"user_id" binding:"required,gt=0"` // 用户ID
	MenuIds []int `json:"menu_ids" binding:"required"`     // 菜单ID列表
}
