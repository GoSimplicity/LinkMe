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

package core

import (
	"context"
)

// 通用 Repository 接口
type GeneralRepository interface {
	Create(ctx context.Context, entity interface{}) error
	FindByID(ctx context.Context, id interface{}) (interface{}, error)
	FindOne(ctx context.Context, query interface{}) (interface{}, error)
	FindAll(ctx context.Context, conditions map[string]interface{}, page, size int) ([]interface{}, error)
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, id interface{}) error
	Transaction(ctx context.Context, fn func() error) error
}

// 通用 Service 接口
type GeneralService interface {
	Create(ctx context.Context, dto interface{}) error
	Get(ctx context.Context, query interface{}) (interface{}, error)
	List(ctx context.Context, query interface{}, page, size int) ([]interface{}, error)
	Update(ctx context.Context, dto interface{}) error
	Delete(ctx context.Context, id interface{}) error
}

// 通用 Database 接口
type GeneralDatabase interface {
	Create(ctx context.Context, entity interface{}) error
	FindOne(ctx context.Context, query interface{}) (interface{}, error)
	FindAll(ctx context.Context, query interface{}, page, size int) ([]interface{}, error)
	Update(ctx context.Context, query interface{}, updates interface{}) error
	Delete(ctx context.Context, query interface{}) error
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
	ExecRaw(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	QueryRaw(ctx context.Context, query string, args ...interface{}) ([]interface{}, error)
}
