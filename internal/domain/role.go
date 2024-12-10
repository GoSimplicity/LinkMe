package domain

// Menu 菜单
type Menu struct {
	ID         int64   `json:"id"`          // 菜单ID
	Name       string  `json:"name"`        // 菜单显示名称
	ParentID   int64   `json:"parent_id"`   // 上级菜单ID,0表示顶级菜单
	Path       string  `json:"path"`        // 前端路由访问路径
	Component  string  `json:"component"`   // 前端组件文件路径
	Icon       string  `json:"icon"`        // 菜单显示图标
	SortOrder  int     `json:"sort_order"`  // 菜单显示顺序,数值越小越靠前
	RouteName  string  `json:"route_name"`  // 前端路由名称,需唯一
	Hidden     int     `json:"hidden"`      // 菜单是否隐藏(0:显示 1:隐藏)
	CreateTime int64   `json:"create_time"` // 记录创建时间戳
	UpdateTime int64   `json:"update_time"` // 记录最后更新时间戳
	IsDeleted  int     `json:"is_deleted"`  // 逻辑删除标记(0:未删除 1:已删除)
	Children   []*Menu `json:"children"`    // 子菜单列表
}

// Api API接口
type Api struct {
	ID          int64  `json:"id"`           // 主键ID
	Name        string `json:"name"`         // API名称
	Path        string `json:"path"`         // API路径
	Method      int    `json:"method"`       // HTTP请求方法(1:GET,2:POST,3:PUT,4:DELETE)
	Description string `json:"description"`  // API描述
	Version     string `json:"version"`      // API版本
	Category    int    `json:"category"`     // API分类(1:系统,2:业务)
	IsPublic    int    `json:"is_public"`    // 是否公开(0:否,1:是)
	CreateTime  int64  `json:"create_time"`  // 创建时间
	UpdateTime  int64  `json:"update_time"`  // 更新时间
	IsDeleted   int    `json:"is_deleted"`   // 是否删除(0:否,1:是)
}

// Role 角色
type Role struct {
	ID          int64  `json:"id"`           // 主键ID
	Name        string `json:"name"`         // 角色名称
	Description string `json:"description"`  // 角色描述
	RoleType    int    `json:"role_type"`    // 角色类型(1:系统角色,2:自定义角色)
	IsDefault   int    `json:"is_default"`   // 是否为默认角色(0:否,1:是)
	CreateTime  int64  `json:"create_time"`  // 创建时间
	UpdateTime  int64  `json:"update_time"`  // 更新时间
	IsDeleted   int    `json:"is_deleted"`   // 是否删除(0:否,1:是)
}
