/*
 * MIT License
 *
 * Copyright (c) 2024 Bamboo
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

package service

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/app/user/repository"
	"github.com/GoSimplicity/LinkMe/internal/core"
	"go.uber.org/zap"
)

type UserService interface {
	core.GeneralService
	Login(ctx context.Context, dto interface{}) (interface{}, error)
	LoginSMS(ctx context.Context, dto interface{}) (interface{}, error)
	SendSMS(ctx context.Context, dto interface{}) error
	SendEmail(ctx context.Context, dto interface{}) error
	RefreshToken(ctx context.Context, dto interface{}) (interface{}, error)
	Logout(ctx context.Context, dto interface{}) error
	GetProfile(ctx context.Context, dto interface{}) (interface{}, error)
	UpdateProfile(ctx context.Context, dto interface{}) error
	ChangePassword(ctx context.Context, dto interface{}) error
	DeleteUser(ctx context.Context, dto interface{}) error
	GetUser(ctx context.Context, dto interface{}) (interface{}, error)
	UpdateUser(ctx context.Context, dto interface{}) error
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

// ChangePassword implements UserService.
func (u *userService) ChangePassword(ctx context.Context, dto interface{}) error {
	panic("unimplemented")
}

// Create implements UserService.
func (u *userService) Create(ctx context.Context, dto interface{}) error {
	panic("unimplemented")
}

// Delete implements UserService.
func (u *userService) Delete(ctx context.Context, id interface{}) error {
	panic("unimplemented")
}

// DeleteUser implements UserService.
func (u *userService) DeleteUser(ctx context.Context, dto interface{}) error {
	panic("unimplemented")
}

// Get implements UserService.
func (u *userService) Get(ctx context.Context, query interface{}) (interface{}, error) {
	panic("unimplemented")
}

// GetProfile implements UserService.
func (u *userService) GetProfile(ctx context.Context, dto interface{}) (interface{}, error) {
	panic("unimplemented")
}

// GetUser implements UserService.
func (u *userService) GetUser(ctx context.Context, dto interface{}) (interface{}, error) {
	panic("unimplemented")
}

// List implements UserService.
func (u *userService) List(ctx context.Context, query interface{}, page int, size int) ([]interface{}, error) {
	panic("unimplemented")
}

// Login implements UserService.
func (u *userService) Login(ctx context.Context, dto interface{}) (interface{}, error) {
	panic("unimplemented")
}

// LoginSMS implements UserService.
func (u *userService) LoginSMS(ctx context.Context, dto interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Logout implements UserService.
func (u *userService) Logout(ctx context.Context, dto interface{}) error {
	panic("unimplemented")
}

// RefreshToken implements UserService.
func (u *userService) RefreshToken(ctx context.Context, dto interface{}) (interface{}, error) {
	panic("unimplemented")
}

// SendEmail implements UserService.
func (u *userService) SendEmail(ctx context.Context, dto interface{}) error {
	panic("unimplemented")
}

// SendSMS implements UserService.
func (u *userService) SendSMS(ctx context.Context, dto interface{}) error {
	panic("unimplemented")
}

// Update implements UserService.
func (u *userService) Update(ctx context.Context, dto interface{}) error {
	panic("unimplemented")
}

// UpdateProfile implements UserService.
func (u *userService) UpdateProfile(ctx context.Context, dto interface{}) error {
	panic("unimplemented")
}

// UpdateUser implements UserService.
func (u *userService) UpdateUser(ctx context.Context, dto interface{}) error {
	panic("unimplemented")
}
