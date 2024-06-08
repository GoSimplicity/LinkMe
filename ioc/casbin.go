package ioc

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"log"
)

// InitCasbin initializes a Casbin enforcer with the given Gorm DB connection and configuration file path.
func InitCasbin(db *gorm.DB, configPath string) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("Failed to create adapter: %v", err)
	}
	enforcer, err := casbin.NewEnforcer(configPath, adapter)
	if err != nil {
		log.Fatalf("Failed to create enforcer: %v", err)
	}
	return enforcer
}
