package repository

import (
	"github.com/GoSimplicity/LinkMe/internal/core"
	"github.com/GoSimplicity/LinkMe/internal/pkg/infra/database/dao"
)

type UserRepository interface {
	core.GeneralRepository
}

type userRepository struct {
	dao dao.UserDao
}

func NewUserRepository(dao dao.UserDao) UserRepository {
	return &userRepository{dao: dao}
}

// Create implements UserRepository.
func (u *userRepository) Create(entity interface{}) error {
	panic("unimplemented")
}

// Delete implements UserRepository.
func (u *userRepository) Delete(id interface{}) error {
	panic("unimplemented")
}

// FindAll implements UserRepository.
func (u *userRepository) FindAll(conditions map[string]interface{}, page int, size int) ([]interface{}, int64, error) {
	panic("unimplemented")
}

// FindByID implements UserRepository.
func (u *userRepository) FindByID(id interface{}) (interface{}, error) {
	panic("unimplemented")
}

// FindOne implements UserRepository.
func (u *userRepository) FindOne(query interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Transaction implements UserRepository.
func (u *userRepository) Transaction(fn func() error) error {
	panic("unimplemented")
}

// Update implements UserRepository.
func (u *userRepository) Update(entity interface{}) error {
	panic("unimplemented")
}
