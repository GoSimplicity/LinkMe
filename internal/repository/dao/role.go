package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	DeletedNo  = 0 // 未删除
	DeletedYes = 1 // 已删除
)

var (
	ErrMenuNotFound = errors.New("菜单不存在")
	ErrInvalidMenu  = errors.New("无效的菜单参数")
)

// Menu 菜单模型
type Menu struct {
	ID         int64   `json:"id" gorm:"primaryKey;column:id;comment:菜单ID"`
	Name       string  `json:"name" gorm:"column:name;type:varchar(50);not null;comment:菜单显示名称"`
	ParentID   int64   `json:"parent_id" gorm:"column:parent_id;default:0;comment:上级菜单ID,0表示顶级菜单"`
	Path       string  `json:"path" gorm:"column:path;type:varchar(255);not null;comment:前端路由访问路径"`
	Component  string  `json:"component" gorm:"column:component;type:varchar(255);not null;comment:前端组件文件路径"`
	Icon       string  `json:"icon" gorm:"column:icon;type:varchar(50);default:'';comment:菜单显示图标"`
	SortOrder  int     `json:"sort_order" gorm:"column:sort_order;default:0;comment:菜单显示顺序,数值越小越靠前"`
	RouteName  string  `json:"route_name" gorm:"column:route_name;type:varchar(50);not null;comment:前端路由名称,需唯一"`
	Hidden     int     `json:"hidden" gorm:"column:hidden;type:tinyint(1);default:0;comment:菜单是否隐藏(0:显示 1:隐藏)"`
	CreateTime int64   `json:"create_time" gorm:"column:create_time;autoCreateTime;comment:记录创建时间戳"`
	UpdateTime int64   `json:"update_time" gorm:"column:update_time;autoUpdateTime;comment:记录最后更新时间戳"`
	IsDeleted  int     `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;comment:逻辑删除标记(0:未删除 1:已删除)"`
	Children   []*Menu `json:"children" gorm:"-"`
}

// Api API模型
type Api struct {
	ID          int64  `json:"id" gorm:"primaryKey;column:id;comment:主键ID"`
	Name        string `json:"name" gorm:"column:name;type:varchar(50);not null;comment:API名称"`
	Path        string `json:"path" gorm:"column:path;type:varchar(255);not null;comment:API路径"`
	Method      int    `json:"method" gorm:"column:method;type:tinyint(1);not null;comment:HTTP请求方法(1:GET,2:POST,3:PUT,4:DELETE)"`
	Description string `json:"description" gorm:"column:description;type:varchar(500);comment:API描述"`
	Version     string `json:"version" gorm:"column:version;type:varchar(20);default:v1;comment:API版本"`
	Category    int    `json:"category" gorm:"column:category;type:tinyint(1);not null;comment:API分类(1:系统,2:业务)"`
	IsPublic    int    `json:"is_public" gorm:"column:is_public;type:tinyint(1);default:0;comment:是否公开(0:否,1:是)"`
	CreateTime  int64  `json:"create_time" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdateTime  int64  `json:"update_time" gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	IsDeleted   int    `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;comment:是否删除(0:否,1:是)"`
}

// Role 角色模型
type Role struct {
	ID          int64  `json:"id" gorm:"primaryKey;column:id;comment:主键ID"`
	Name        string `json:"name" gorm:"column:name;type:varchar(50);not null;unique;comment:角色名称"`
	Description string `json:"description" gorm:"column:description;type:varchar(255);comment:角色描述"`
	RoleType    int    `json:"role_type" gorm:"column:role_type;type:tinyint(1);not null;comment:角色类型(1:系统角色,2:自定义角色)"`
	IsDefault   int    `json:"is_default" gorm:"column:is_default;type:tinyint(1);default:0;comment:是否为默认角色(0:否,1:是)"`
	CreateTime  int64  `json:"create_time" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdateTime  int64  `json:"update_time" gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	IsDeleted   int    `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;comment:是否删除(0:否,1:是)"`
}

// PermissionDAO 定义了权限数据库访问接口
type PermissionDAO interface {
	// 菜单管理
	CreateMenu(ctx context.Context, menu *Menu) error
	GetMenuById(ctx context.Context, id int) (*Menu, error)
	UpdateMenu(ctx context.Context, menu *Menu) error
	DeleteMenu(ctx context.Context, id int) error
	ListMenus(ctx context.Context, page, pageSize int) ([]*Menu, int, error)
	GetMenuTree(ctx context.Context) ([]*Menu, error)

	// API接口管理
	CreateApi(ctx context.Context, api *Api) error
	GetApiById(ctx context.Context, id int) (*Api, error)
	UpdateApi(ctx context.Context, api *Api) error
	DeleteApi(ctx context.Context, id int) error
	ListApis(ctx context.Context, page, pageSize int) ([]*Api, int, error)

	// 角色管理
	CreateRole(ctx context.Context, role *Role) error
	GetRoleById(ctx context.Context, id int) (*Role, error)
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, id int) error
	ListRoles(ctx context.Context, page, pageSize int) ([]*Role, int, error)
	AssignPermissions(ctx context.Context, roleId int, menuIds []int, apiIds []int) error
	AssignRoleToUser(ctx context.Context, userId int, roleIds []int) error
	AssignApiPermissionsToUser(ctx context.Context, userId int, apiIds []int) error
	AssignMenuPermissionsToUser(ctx context.Context, userId int, menuIds []int) error
	RemoveUserApiPermissions(ctx context.Context, userId int, apiIds []int) error
	RemoveUserMenuPermissions(ctx context.Context, userId int, menuIds []int) error
	RemoveRoleFromUser(ctx context.Context, userId int, roleIds []int) error
	RemoveRoleApiPermissions(ctx context.Context, roleIds []int, apiIds []int) error
	RemoveRoleMenuPermissions(ctx context.Context, roleIds []int, menuIds []int) error
}

// permissionDAO 是 PermissionDAO 的实现
type permissionDAO struct {
	db       *gorm.DB
	l        *zap.Logger
	enforcer *casbin.Enforcer
}

// NewPermissionDAO 创建一个新的 PermissionDAO
func NewPermissionDAO(db *gorm.DB, l *zap.Logger, enforcer *casbin.Enforcer) PermissionDAO {
	return &permissionDAO{
		db:       db,
		l:        l,
		enforcer: enforcer,
	}
}

func (p permissionDAO) CreateMenu(ctx context.Context, menu *Menu) error {
	if menu == nil {
		return ErrInvalidMenu
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查父菜单是否存在
		if menu.ParentID != 0 {
			var count int64
			if err := tx.Model(&Menu{}).Where("id = ? AND is_deleted = ?", menu.ParentID, DeletedNo).Count(&count).Error; err != nil {
				return fmt.Errorf("检查父菜单失败: %v", err)
			}
			if count == 0 {
				return errors.New("父菜单不存在")
			}
		}

		now := time.Now().Unix()
		menu.CreateTime = now
		menu.UpdateTime = now

		return tx.Create(menu).Error
	})
}

func (p permissionDAO) GetMenuById(ctx context.Context, id int) (*Menu, error) {
	if id <= 0 {
		return nil, errors.New("无效的菜单ID")
	}

	var menu Menu
	if err := p.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, DeletedNo).First(&menu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMenuNotFound
		}
		return nil, fmt.Errorf("查询菜单失败: %v", err)
	}

	return &menu, nil
}

func (p permissionDAO) UpdateMenu(ctx context.Context, menu *Menu) error {
	if menu == nil {
		return errors.New("菜单对象不能为空")
	}
	if menu.ID <= 0 {
		return errors.New("无效的菜单ID")
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查父菜单是否存在且不能将菜单设置为自己的子菜单
		if menu.ParentID != 0 {
			if menu.ParentID == menu.ID {
				return errors.New("不能将菜单设置为自己的子菜单")
			}
			var count int64
			if err := tx.Model(&Menu{}).Where("id = ? AND is_deleted = ?", menu.ParentID, DeletedNo).Count(&count).Error; err != nil {
				return fmt.Errorf("检查父菜单失败: %v", err)
			}
			if count == 0 {
				return errors.New("父菜单不存在")
			}
		}

		updates := map[string]interface{}{
			"name":        menu.Name,
			"parent_id":   menu.ParentID,
			"path":        menu.Path,
			"component":   menu.Component,
			"icon":        menu.Icon,
			"sort_order":  menu.SortOrder,
			"route_name":  menu.RouteName,
			"hidden":      menu.Hidden,
			"update_time": time.Now().Unix(),
		}

		result := tx.Model(&Menu{}).
			Where("id = ? AND is_deleted = ?", menu.ID, DeletedNo).
			Updates(updates)
		if result.Error != nil {
			return fmt.Errorf("更新菜单失败: %v", result.Error)
		}
		if result.RowsAffected == 0 {
			return errors.New("菜单不存在或已被删除")
		}

		return nil
	})
}

func (p permissionDAO) DeleteMenu(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("无效的菜单ID")
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查是否有子菜单
		var count int64
		if err := tx.Model(&Menu{}).Where("parent_id = ? AND is_deleted = ?", id, DeletedNo).Count(&count).Error; err != nil {
			return fmt.Errorf("检查子菜单失败: %v", err)
		}
		if count > 0 {
			return errors.New("存在子菜单,不能删除")
		}

		updates := map[string]interface{}{
			"is_deleted":  DeletedYes,
			"update_time": time.Now().Unix(),
		}
		result := tx.Model(&Menu{}).Where("id = ? AND is_deleted = ?", id, DeletedNo).Updates(updates)
		if result.Error != nil {
			return fmt.Errorf("删除菜单失败: %v", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrMenuNotFound
		}
		return nil
	})
}

func (p permissionDAO) ListMenus(ctx context.Context, page, pageSize int) ([]*Menu, int, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("无效的分页参数")
	}

	var menus []*Menu
	var total int64

	db := p.db.WithContext(ctx).Model(&Menu{}).Where("is_deleted = ?", DeletedNo)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取菜单总数失败: %v", err)
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("sort_order ASC, id DESC").Find(&menus).Error; err != nil {
		return nil, 0, fmt.Errorf("查询菜单列表失败: %v", err)
	}

	return menus, int(total), nil
}

func (p permissionDAO) GetMenuTree(ctx context.Context) ([]*Menu, error) {
	// 预分配合适的初始容量
	menus := make([]*Menu, 0, 50)

	// 使用索引字段优化查询
	if err := p.db.WithContext(ctx).
		Select("id, name, parent_id, path, component, icon, sort_order, route_name, hidden, create_time, update_time").
		Where("is_deleted = ?", DeletedNo).
		Order("sort_order ASC, id ASC").
		Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("查询菜单列表失败: %v", err)
	}

	// 预分配map容量
	menuMap := make(map[int64]*Menu, len(menus))
	rootMenus := make([]*Menu, 0, len(menus)/3)

	// 第一次遍历,建立ID到菜单的映射
	for _, menu := range menus {
		menu.Children = make([]*Menu, 0, 4)
		menuMap[menu.ID] = menu
	}

	// 第二次遍历,构建树形结构
	for _, menu := range menus {
		if menu.ParentID == 0 {
			rootMenus = append(rootMenus, menu)
		} else {
			if parent, exists := menuMap[menu.ParentID]; exists {
				parent.Children = append(parent.Children, menu)
			} else {
				// 如果找不到父节点,作为根节点处理
				rootMenus = append(rootMenus, menu)
			}
		}
	}

	return rootMenus, nil
}

func (p permissionDAO) CreateApi(ctx context.Context, api *Api) error {
	if api == nil {
		return gorm.ErrRecordNotFound
	}

	api.CreateTime = time.Now().Unix()
	api.UpdateTime = time.Now().Unix()

	return p.db.WithContext(ctx).Create(api).Error
}

func (p permissionDAO) GetApiById(ctx context.Context, id int) (*Api, error) {
	var api Api

	if err := p.db.WithContext(ctx).Where("id = ? AND is_deleted = 0", id).First(&api).Error; err != nil {
		return nil, err
	}

	return &api, nil
}

func (p permissionDAO) UpdateApi(ctx context.Context, api *Api) error {
	if api == nil {
		return gorm.ErrRecordNotFound
	}

	updates := map[string]interface{}{
		"name":        api.Name,
		"path":        api.Path,
		"method":      api.Method,
		"description": api.Description,
		"version":     api.Version,
		"category":    api.Category,
		"is_public":   api.IsPublic,
		"update_time": time.Now().Unix(),
	}

	return p.db.WithContext(ctx).
		Model(&Api{}).
		Where("id = ? AND is_deleted = 0", api.ID).
		Updates(updates).Error
}

func (p permissionDAO) DeleteApi(ctx context.Context, id int) error {
	updates := map[string]interface{}{
		"is_deleted":  1,
		"update_time": time.Now().Unix(),
	}

	return p.db.WithContext(ctx).Model(&Api{}).Where("id = ? AND is_deleted = 0", id).Updates(updates).Error
}

func (p permissionDAO) ListApis(ctx context.Context, page, pageSize int) ([]*Api, int, error) {
	var apis []*Api
	var total int64

	db := p.db.WithContext(ctx).Model(&Api{}).Where("is_deleted = 0")

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("id ASC").Find(&apis).Error; err != nil {
		return nil, 0, err
	}

	return apis, int(total), nil
}

func (p permissionDAO) CreateRole(ctx context.Context, role *Role) error {
	if role == nil {
		return errors.New("角色对象不能为空")
	}

	if role.Name == "" {
		return errors.New("角色名称不能为空")
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查角色名是否已存在
		var count int64
		if err := tx.Model(&Role{}).Where("name = ? AND is_deleted = ?", role.Name, 0).Count(&count).Error; err != nil {
			return fmt.Errorf("检查角色名称失败: %v", err)
		}
		if count > 0 {
			return errors.New("角色名称已存在")
		}

		// 设置创建时间和更新时间
		now := time.Now().Unix()
		role.CreateTime = now
		role.UpdateTime = now
		role.IsDeleted = 0

		// 创建角色
		if err := tx.Create(role).Error; err != nil {
			return fmt.Errorf("创建角色失败: %v", err)
		}

		return nil
	})
}

func (p permissionDAO) GetRoleById(ctx context.Context, id int) (*Role, error) {
	if id <= 0 {
		return nil, errors.New("无效的角色ID")
	}

	var role Role
	if err := p.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, 0).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询角色失败: %v", err)
	}

	return &role, nil
}

func (p permissionDAO) UpdateRole(ctx context.Context, role *Role) error {
	if role == nil {
		return errors.New("角色对象不能为空")
	}
	if role.ID <= 0 {
		return errors.New("无效的角色ID")
	}
	if role.Name == "" {
		return errors.New("角色名称不能为空")
	}

	// 检查角色名是否已被其他角色使用
	var count int64
	if err := p.db.WithContext(ctx).Model(&Role{}).
		Where("name = ? AND id != ? AND is_deleted = ?", role.Name, role.ID, 0).
		Count(&count).Error; err != nil {
		return fmt.Errorf("检查角色名称失败: %v", err)
	}
	if count > 0 {
		return errors.New("角色名称已被使用")
	}

	updates := map[string]interface{}{
		"name":        role.Name,
		"description": role.Description,
		"role_type":   role.RoleType,
		"is_default":  role.IsDefault,
		"update_time": time.Now().Unix(),
	}

	result := p.db.WithContext(ctx).
		Model(&Role{}).
		Where("id = ? AND is_deleted = ?", role.ID, 0).
		Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新角色失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("角色不存在或已被删除")
	}

	return nil
}

func (p permissionDAO) DeleteRole(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("无效的角色ID")
	}

	// 检查是否为默认角色
	var role Role
	if err := p.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, 0).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %v", err)
	}

	if role.IsDefault == 1 {
		return errors.New("默认角色不能删除")
	}

	updates := map[string]interface{}{
		"is_deleted":  1,
		"update_time": time.Now().Unix(),
	}

	result := p.db.WithContext(ctx).Model(&Role{}).Where("id = ? AND is_deleted = ?", id, 0).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("删除角色失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("角色不存在或已被删除")
	}

	// 删除角色关联的权限
	if _, err := p.enforcer.DeleteRole(role.Name); err != nil {
		return fmt.Errorf("删除角色权限失败: %v", err)
	}

	return nil
}

func (p permissionDAO) ListRoles(ctx context.Context, page, pageSize int) ([]*Role, int, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("无效的分页参数")
	}

	var roles []*Role
	var total int64

	db := p.db.WithContext(ctx).Model(&Role{}).Where("is_deleted = ?", 0)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取角色总数失败: %v", err)
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("id ASC").Find(&roles).Error; err != nil {
		return nil, 0, fmt.Errorf("获取角色列表失败: %v", err)
	}

	return roles, int(total), nil
}

func (p permissionDAO) AssignPermissions(ctx context.Context, roleId int, menuIds []int, apiIds []int) error {
	const batchSize = 1000

	if roleId <= 0 {
		return errors.New("无效的角色ID")
	}

	// 检查角色是否存在
	role, err := p.GetRoleById(ctx, roleId)
	if err != nil {
		return fmt.Errorf("获取角色失败: %v", err)
	}
	if role == nil {
		return errors.New("角色不存在")
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除原有的casbin规则
		if _, err := p.enforcer.DeleteRolesForUser(role.Name); err != nil {
			return fmt.Errorf("删除原有权限失败: %v", err)
		}

		// 添加菜单权限
		if err := p.assignMenuPermissions(role.Name, menuIds, batchSize); err != nil {
			return err
		}

		// 添加API权限
		if err := p.assignAPIPermissions(ctx, role.Name, apiIds, batchSize); err != nil {
			return err
		}

		// 加载最新的策略
		if err := p.enforcer.LoadPolicy(); err != nil {
			return fmt.Errorf("加载策略失败: %v", err)
		}

		return nil
	})
}

// AssignApiPermissionsToUser 添加用户api权限
func (p *permissionDAO) AssignApiPermissionsToUser(ctx context.Context, userId int, apiIds []int) error {
	if userId <= 0 {
		return errors.New("无效的用户ID")
	}

	if len(apiIds) == 0 {
		return nil
	}

	// 查询API信息
	var apis []*Api
	if err := p.db.WithContext(ctx).Where("id IN ? AND is_deleted = ?", apiIds, DeletedNo).Find(&apis).Error; err != nil {
		return fmt.Errorf("查询API失败: %v", err)
	}

	if len(apis) == 0 {
		return errors.New("未找到有效的API")
	}

	userStr := fmt.Sprintf("%d", userId)

	// 构建需要添加的策略
	var policies [][]string
	methodMap := map[int]string{1: "GET", 2: "POST", 3: "PUT", 4: "DELETE", 5: "PATCH", 6: "OPTIONS", 7: "HEAD"}

	for _, api := range apis {
		method, ok := methodMap[api.Method]
		if !ok {
			return fmt.Errorf("无效的HTTP方法: %d", api.Method)
		}
		policies = append(policies, []string{userStr, api.Path, method})
	}

	// 批量添加策略
	if _, err := p.enforcer.AddPolicies(policies); err != nil {
		return fmt.Errorf("添加权限策略失败: %v", err)
	}

	// 重新加载策略
	if err := p.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载策略失败: %v", err)
	}

	return nil
}

// AssignMenuPermissionsToUser 添加用户菜单权限
func (p *permissionDAO) AssignMenuPermissionsToUser(ctx context.Context, userId int, menuIds []int) error {
	if userId <= 0 {
		return errors.New("无效的用户ID")
	}

	if len(menuIds) == 0 {
		return nil
	}

	// 查询菜单信息
	var menus []*Menu
	if err := p.db.WithContext(ctx).Where("id IN ? AND is_deleted = ?", menuIds, DeletedNo).Find(&menus).Error; err != nil {
		return fmt.Errorf("查询菜单失败: %v", err)
	}

	if len(menus) == 0 {
		return errors.New("未找到有效的菜单")
	}

	userStr := fmt.Sprintf("%d", userId)

	// 构建需要添加的策略
	var policies [][]string
	for _, menu := range menus {
		policies = append(policies, []string{userStr, fmt.Sprintf("menu:%d", menu.ID), "read"})
	}

	// 批量添加策略
	if _, err := p.enforcer.AddPolicies(policies); err != nil {
		return fmt.Errorf("添加权限策略失败: %v", err)
	}

	// 重新加载策略
	if err := p.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载策略失败: %v", err)
	}

	return nil
}

func (p permissionDAO) AssignRoleToUser(ctx context.Context, userId int, roleIds []int) error {
	if userId <= 0 {
		return errors.New("无效的用户ID")
	}

	if len(roleIds) == 0 {
		return nil
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 获取角色信息
		var roles []*Role
		if err := tx.Where("id IN ? AND is_deleted = ?", roleIds, 0).Find(&roles).Error; err != nil {
			return fmt.Errorf("获取角色信息失败: %v", err)
		}

		if len(roles) == 0 {
			return errors.New("未找到有效的角色")
		}

		// 构建casbin规则
		policies := make([][]string, 0, len(roles))
		for _, role := range roles {
			policies = append(policies, []string{fmt.Sprintf("%d", userId), role.Name})
		}

		// 添加用户角色关联
		if _, err := p.enforcer.AddGroupingPolicies(policies); err != nil {
			return fmt.Errorf("添加用户角色关联失败: %v", err)
		}

		// 加载最新的策略
		if err := p.enforcer.LoadPolicy(); err != nil {
			return fmt.Errorf("加载策略失败: %v", err)
		}

		return nil
	})
}

// RemoveRoleFromUser 移除用户角色
func (p permissionDAO) RemoveRoleFromUser(ctx context.Context, userId int, roleIds []int) error {
	if userId <= 0 {
		return errors.New("无效的用户ID")
	}

	if len(roleIds) == 0 {
		return nil
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 获取角色信息
		var roles []*Role
		if err := tx.Where("id IN ? AND is_deleted = ?", roleIds, DeletedNo).Find(&roles).Error; err != nil {
			return fmt.Errorf("获取角色信息失败: %v", err)
		}

		if len(roles) == 0 {
			return errors.New("未找到有效的角色")
		}

		// 构建需要移除的规则
		policies := make([][]string, 0, len(roles))
		for _, role := range roles {
			// 检查用户是否拥有该角色
			hasRole, err := p.enforcer.HasGroupingPolicy(fmt.Sprintf("%d", userId), role.Name)
			if err != nil {
				return fmt.Errorf("检查用户角色失败: %v", err)
			}
			if hasRole {
				policies = append(policies, []string{fmt.Sprintf("%d", userId), role.Name})
			}
		}

		// 如果没有需要移除的规则,直接返回
		if len(policies) == 0 {
			return nil
		}

		// 移除用户角色关联
		if _, err := p.enforcer.RemoveGroupingPolicies(policies); err != nil {
			return fmt.Errorf("移除用户角色关联失败: %v", err)
		}

		// 加载最新的策略
		if err := p.enforcer.LoadPolicy(); err != nil {
			return fmt.Errorf("加载策略失败: %v", err)
		}

		return nil
	})
}

// RemoveUserApiPermissions 移除用户api权限
func (p *permissionDAO) RemoveUserApiPermissions(ctx context.Context, userId int, apiIds []int) error {
	if userId <= 0 {
		return errors.New("无效的用户ID")
	}

	if len(apiIds) == 0 {
		return nil
	}

	// 查询API信息
	var apis []*Api
	if err := p.db.WithContext(ctx).Where("id IN ? AND is_deleted = ?", apiIds, DeletedNo).Find(&apis).Error; err != nil {
		return fmt.Errorf("查询API失败: %v", err)
	}

	if len(apis) == 0 {
		return errors.New("未找到有效的API")
	}

	// 构建需要移除的策略
	var policies [][]string
	methodMap := map[int]string{1: "GET", 2: "POST", 3: "PUT", 4: "DELETE", 5: "PATCH", 6: "OPTIONS", 7: "HEAD"}

	for _, api := range apis {
		method, ok := methodMap[api.Method]
		if !ok {
			return fmt.Errorf("无效的HTTP方法: %d", api.Method)
		}
		policies = append(policies, []string{fmt.Sprintf("%d", userId), api.Path, method})
	}

	// 批量移除策略
	if _, err := p.enforcer.RemovePolicies(policies); err != nil {
		return fmt.Errorf("移除权限策略失败: %v", err)
	}

	// 重新加载策略
	if err := p.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载策略失败: %v", err)
	}

	return nil
}

// RemoveUserMenuPermissions 移除用户菜单权限
func (p *permissionDAO) RemoveUserMenuPermissions(ctx context.Context, userId int, menuIds []int) error {
	if userId <= 0 {
		return errors.New("无效的用户ID")
	}

	if len(menuIds) == 0 {
		return nil
	}

	// 查询菜单信息
	var menus []*Menu
	if err := p.db.WithContext(ctx).Where("id IN ? AND is_deleted = ?", menuIds, DeletedNo).Find(&menus).Error; err != nil {
		return fmt.Errorf("查询菜单失败: %v", err)
	}

	if len(menus) == 0 {
		return errors.New("未找到有效的菜单")
	}

	// 构建需要移除的策略
	var policies [][]string
	for _, menu := range menus {
		policies = append(policies, []string{fmt.Sprintf("%d", userId), fmt.Sprintf("menu:%d", menu.ID), "read"})
	}

	// 批量移除策略
	if _, err := p.enforcer.RemovePolicies(policies); err != nil {
		return fmt.Errorf("移除权限策略失败: %v", err)
	}

	// 重新加载策略
	if err := p.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载策略失败: %v", err)
	}

	return nil
}

// RemoveRoleApiPermissions 批量移除角色对应api权限
func (p *permissionDAO) RemoveRoleApiPermissions(ctx context.Context, roleIds []int, apiIds []int) error {
	if len(roleIds) == 0 || len(apiIds) == 0 {
		return nil
	}

	// 查询角色名称
	var roles []*Role
	if err := p.db.WithContext(ctx).Where("id IN ? AND is_deleted = ?", roleIds, DeletedNo).Find(&roles).Error; err != nil {
		return fmt.Errorf("查询角色失败: %v", err)
	}

	if len(roles) == 0 {
		return errors.New("未找到有效的角色")
	}

	// 查询API信息
	var apis []*Api
	if err := p.db.WithContext(ctx).Where("id IN ? AND is_deleted = ?", apiIds, DeletedNo).Find(&apis).Error; err != nil {
		return fmt.Errorf("查询API失败: %v", err)
	}

	if len(apis) == 0 {
		return errors.New("未找到有效的API")
	}

	// 构建需要移除的策略
	var policies [][]string
	methodMap := map[int]string{1: "GET", 2: "POST", 3: "PUT", 4: "DELETE"}

	for _, role := range roles {
		for _, api := range apis {
			method, ok := methodMap[api.Method]
			if !ok {
				return fmt.Errorf("无效的HTTP方法: %d", api.Method)
			}
			policies = append(policies, []string{role.Name, api.Path, method})
		}
	}

	// 批量移除策略
	if _, err := p.enforcer.RemovePolicies(policies); err != nil {
		return fmt.Errorf("移除权限策略失败: %v", err)
	}

	// 重新加载策略
	if err := p.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载策略失败: %v", err)
	}

	return nil
}

// RemoveRoleMenuPermissions 批量移除角色对应菜单权限
func (p *permissionDAO) RemoveRoleMenuPermissions(ctx context.Context, roleIds []int, menuIds []int) error {
	if len(roleIds) == 0 || len(menuIds) == 0 {
		return nil
	}

	// 查询角色名称
	var roles []*Role
	if err := p.db.WithContext(ctx).Where("id IN ? AND is_deleted = ?", roleIds, DeletedNo).Find(&roles).Error; err != nil {
		return fmt.Errorf("查询角色失败: %v", err)
	}

	if len(roles) == 0 {
		return errors.New("未找到有效的角色")
	}

	// 查询菜单信息
	var menus []*Menu
	if err := p.db.WithContext(ctx).Where("id IN ? AND is_deleted = ?", menuIds, DeletedNo).Find(&menus).Error; err != nil {
		return fmt.Errorf("查询菜单失败: %v", err)
	}

	if len(menus) == 0 {
		return errors.New("未找到有效的菜单")
	}

	// 构建需要移除的策略
	var policies [][]string
	for _, role := range roles {
		for _, menu := range menus {
			policies = append(policies, []string{role.Name, fmt.Sprintf("menu:%d", menu.ID), "read"})
		}
	}

	// 批量移除策略
	if _, err := p.enforcer.RemovePolicies(policies); err != nil {
		return fmt.Errorf("移除权限策略失败: %v", err)
	}

	// 重新加载策略
	if err := p.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载策略失败: %v", err)
	}

	return nil
}

// assignMenuPermissions 分配菜单权限
func (p permissionDAO) assignMenuPermissions(roleName string, menuIds []int, batchSize int) error {
	if roleName == "" {
		return errors.New("角色名称不能为空")
	}

	// 如果菜单ID列表为空,直接返回
	if len(menuIds) == 0 {
		return nil
	}

	// 构建casbin策略规则
	policies := make([][]string, 0, len(menuIds))
	for _, menuId := range menuIds {
		if menuId <= 0 {
			return fmt.Errorf("无效的菜单ID: %d", menuId)
		}
		policies = append(policies, []string{roleName, fmt.Sprintf("menu:%d", menuId), "read"})
	}

	// 批量添加策略
	return p.batchAddPolicies(policies, batchSize)
}

// assignAPIPermissions 分配API权限
func (p permissionDAO) assignAPIPermissions(ctx context.Context, roleName string, apiIds []int, batchSize int) error {
	if roleName == "" {
		return errors.New("角色名称不能为空")
	}

	// 如果API ID列表为空,直接返回
	if len(apiIds) == 0 {
		return nil
	}

	// HTTP方法映射表
	methodMap := map[int]string{
		1: "GET",
		2: "POST",
		3: "PUT",
		4: "DELETE",
		5: "PATCH",
		6: "OPTIONS",
		7: "HEAD",
	}

	// 构建casbin策略规则
	policies := make([][]string, 0, len(apiIds))
	for _, apiId := range apiIds {
		if apiId <= 0 {
			return fmt.Errorf("无效的API ID: %d", apiId)
		}

		// 获取API信息
		api, err := p.GetApiById(ctx, apiId)
		if err != nil {
			return fmt.Errorf("获取API信息失败: %v", err)
		}

		if api == nil {
			return fmt.Errorf("API不存在: %d", apiId)
		}

		// 获取HTTP方法
		method, ok := methodMap[api.Method]
		if !ok {
			return fmt.Errorf("无效的HTTP方法: %d", api.Method)
		}

		policies = append(policies, []string{roleName, api.Path, method})
	}

	// 批量添加策略
	return p.batchAddPolicies(policies, batchSize)
}

// batchAddPolicies 批量添加策略
func (p permissionDAO) batchAddPolicies(policies [][]string, batchSize int) error {
	if len(policies) == 0 {
		return nil
	}

	if batchSize <= 0 {
		return errors.New("无效的批次大小")
	}

	// 按批次处理策略规则
	for i := 0; i < len(policies); i += batchSize {
		end := i + batchSize
		if end > len(policies) {
			end = len(policies)
		}

		// 添加一批策略规则
		if _, err := p.enforcer.AddPolicies(policies[i:end]); err != nil {
			return fmt.Errorf("添加权限策略失败: %v", err)
		}
	}

	// 加载最新的策略
	if err := p.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载策略失败: %v", err)
	}

	return nil
}
