package dao

import (
	"github.com/GoSimplicity/LinkMe/internal/core"
	"gorm.io/gorm"
)

type UserDao interface {
	core.GeneralDatabase
}

type userDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &userDao{db: db}
}

// Create implements UserDao.
func (u *userDao) Create(ctx interface{}, entity interface{}) error {
	panic("unimplemented")
}

// Delete implements UserDao.
func (u *userDao) Delete(ctx interface{}, query interface{}) error {
	panic("unimplemented")
}

// ExecRaw implements UserDao.
func (u *userDao) ExecRaw(query string, args ...interface{}) (interface{}, error) {
	panic("unimplemented")
}

// FindAll implements UserDao.
func (u *userDao) FindAll(ctx interface{}, query interface{}, page int, size int) ([]interface{}, int64, error) {
	panic("unimplemented")
}

// FindOne implements UserDao.
func (u *userDao) FindOne(ctx interface{}, query interface{}) (interface{}, error) {
	panic("unimplemented")
}

// QueryRaw implements UserDao.
func (u *userDao) QueryRaw(query string, args ...interface{}) ([]interface{}, error) {
	panic("unimplemented")
}

// Transaction implements UserDao.
func (u *userDao) Transaction(ctx interface{}, fn func(ctx interface{}) error) error {
	panic("unimplemented")
}

// Update implements UserDao.
func (u *userDao) Update(ctx interface{}, query interface{}, updates interface{}) error {
	panic("unimplemented")
}
