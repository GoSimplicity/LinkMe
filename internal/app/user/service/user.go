package service

import "github.com/GoSimplicity/LinkMe/internal/core"

type UserService interface {
	core.GeneralService
	//TODO: 添加用户服务方法
}

type userService struct {
}

func NewUserService() UserService {
	return &userService{}
}

// Create implements UserService.
func (u *userService) Create(ctx interface{}, dto interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Delete implements UserService.
func (u *userService) Delete(ctx interface{}, id interface{}) error {
	panic("unimplemented")
}

// Get implements UserService.
func (u *userService) Get(ctx interface{}, query interface{}) (interface{}, error) {
	panic("unimplemented")
}

// List implements UserService.
func (u *userService) List(ctx interface{}, query interface{}, page int, size int) ([]interface{}, int64, error) {
	panic("unimplemented")
}

// Update implements UserService.
func (u *userService) Update(ctx interface{}, dto interface{}) error {
	panic("unimplemented")
}
