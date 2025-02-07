package ioc

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"log"
)

// InitCasbin 初始化casbin权限管理器
func InitCasbin(db *gorm.DB) *casbin.Enforcer {
	// 创建gorm适配器,用于将权限规则存储到数据库中
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("创建适配器失败: %v", err)
	}

	// 创建enforcer实例,使用配置文件中的模型定义和数据库适配器
	enforcer, err := casbin.NewEnforcer("config/model.conf", adapter)
	if err != nil {
		log.Fatalf("创建enforcer失败: %v", err)
	}
	return enforcer
}
