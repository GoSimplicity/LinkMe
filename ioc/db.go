package ioc

import (
	"LinkMe/internal/dao"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type config struct {
	DSN string `yaml:"dsn"`
}

func InitDB() *gorm.DB {
	var c config

	if err := viper.UnmarshalKey("db", &c); err != nil {
		panic(fmt.Errorf("初始化失败：%v", err))
	}
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err = dao.InitTables(db); err != nil {
		panic(err)
	}
	return db
}
