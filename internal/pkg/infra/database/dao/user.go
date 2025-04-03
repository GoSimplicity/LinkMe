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

package dao

import (
	"context"

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
func (u *userDao) Create(ctx context.Context, entity interface{}) error {
	panic("unimplemented")
}

// Delete implements UserDao.
func (u *userDao) Delete(ctx context.Context, query interface{}) error {
	panic("unimplemented")
}

// ExecRaw implements UserDao.
func (u *userDao) ExecRaw(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	panic("unimplemented")
}

// FindAll implements UserDao.
func (u *userDao) FindAll(ctx context.Context, query interface{}, page int, size int) ([]interface{}, error) {
	panic("unimplemented")
}

// FindOne implements UserDao.
func (u *userDao) FindOne(ctx context.Context, query interface{}) (interface{}, error) {
	panic("unimplemented")
}

// QueryRaw implements UserDao.
func (u *userDao) QueryRaw(ctx context.Context, query string, args ...interface{}) ([]interface{}, error) {
	panic("unimplemented")
}

// Transaction implements UserDao.
func (u *userDao) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	panic("unimplemented")
}

// Update implements UserDao.
func (u *userDao) Update(ctx context.Context, query interface{}, updates interface{}) error {
	panic("unimplemented")
}
