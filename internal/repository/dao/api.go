package dao

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Api struct {
	ID          int    `json:"id" gorm:"primaryKey;column:id;comment:主键ID"`
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

type ApiDAO interface {
	CreateApi(ctx context.Context, api *Api) error
	GetApiById(ctx context.Context, id int) (*Api, error)
	UpdateApi(ctx context.Context, api *Api) error
	DeleteApi(ctx context.Context, id int) error
	ListApis(ctx context.Context, page, pageSize int) ([]*Api, int, error)
}

type apiDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

func NewApiDAO(db *gorm.DB, l *zap.Logger) ApiDAO {
	return &apiDAO{
		db: db,
		l:  l,
	}
}

// CreateApi 创建新的API记录
func (a *apiDAO) CreateApi(ctx context.Context, api *Api) error {
	if api == nil {
		return gorm.ErrRecordNotFound
	}

	api.CreateTime = time.Now().Unix()
	api.UpdateTime = time.Now().Unix()

	return a.db.WithContext(ctx).Create(api).Error
}

// GetApiById 根据ID获取API记录
func (a *apiDAO) GetApiById(ctx context.Context, id int) (*Api, error) {
	var api Api

	if err := a.db.WithContext(ctx).Where("id = ? AND is_deleted = 0", id).First(&api).Error; err != nil {
		return nil, err
	}

	return &api, nil
}

// UpdateApi 更新API记录
func (a *apiDAO) UpdateApi(ctx context.Context, api *Api) error {
	if api == nil {
		return gorm.ErrRecordNotFound
	}

	// 确保API路径以api:开头
	if !strings.HasPrefix(api.Path, "api:") {
		api.Path = "api:" + api.Path
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

	return a.db.WithContext(ctx).
		Model(&Api{}).
		Where("id = ? AND is_deleted = 0", api.ID).
		Updates(updates).Error
}

// DeleteApi 软删除API记录
func (a *apiDAO) DeleteApi(ctx context.Context, id int) error {
	updates := map[string]interface{}{
		"is_deleted":  1,
		"update_time": time.Now().Unix(),
	}

	return a.db.WithContext(ctx).Model(&Api{}).Where("id = ? AND is_deleted = 0", id).Updates(updates).Error
}

// ListApis 分页获取API列表
func (a *apiDAO) ListApis(ctx context.Context, page, pageSize int) ([]*Api, int, error) {
	var apis []*Api
	var total int64

	// 构建基础查询
	db := a.db.WithContext(ctx).Model(&Api{}).Where("is_deleted = 0")

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
