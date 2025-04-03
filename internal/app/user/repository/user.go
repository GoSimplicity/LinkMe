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

package repository

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/core"
	"github.com/GoSimplicity/LinkMe/internal/pkg/infra/database/dao"
)

type UserRepository interface {
	core.GeneralRepository
}

type userRepository struct {
	dao dao.UserDao
}

// Create implements UserRepository.
func (u *userRepository) Create(ctx context.Context, entity interface{}) error {
	panic("unimplemented")
}

// Delete implements UserRepository.
func (u *userRepository) Delete(ctx context.Context, id interface{}) error {
	panic("unimplemented")
}

// FindAll implements UserRepository.
func (u *userRepository) FindAll(ctx context.Context, conditions map[string]interface{}, page int, size int) ([]interface{}, error) {
	panic("unimplemented")
}

// FindByID implements UserRepository.
func (u *userRepository) FindByID(ctx context.Context, id interface{}) (interface{}, error) {
	panic("unimplemented")
}

// FindOne implements UserRepository.
func (u *userRepository) FindOne(ctx context.Context, query interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Transaction implements UserRepository.
func (u *userRepository) Transaction(ctx context.Context, fn func() error) error {
	panic("unimplemented")
}

// Update implements UserRepository.
func (u *userRepository) Update(ctx context.Context, entity interface{}) error {
	panic("unimplemented")
}

func NewUserRepository(dao dao.UserDao) UserRepository {
	return &userRepository{dao: dao}
}
