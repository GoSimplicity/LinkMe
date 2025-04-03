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
	"fmt"
	"net/http"
)

type ErrorCode int

const (
	// 成功
	Success ErrorCode = 0

	// 通用错误 (1-999)
	UnknownError     ErrorCode = 1   // 未知错误
	InternalError    ErrorCode = 2   // 内部错误
	InvalidParameter ErrorCode = 100 // 参数错误
	Unauthorized     ErrorCode = 401 // 未授权
	Forbidden        ErrorCode = 403 // 禁止访问
	NotFound         ErrorCode = 404 // 资源不存在
	Timeout          ErrorCode = 408 // 请求超时

	// 业务错误 (1000-1999)
	BusinessError ErrorCode = 1000 // 业务通用错误

	// 数据库错误 (2000-2999)
	DBError           ErrorCode = 2000 // 数据库通用错误
	DBConnectionError ErrorCode = 2001 // 数据库连接错误
	DBQueryError      ErrorCode = 2002 // 数据库查询错误

	// 缓存错误 (3000-3999)
	CacheError ErrorCode = 3000 // 缓存通用错误

	// 第三方服务错误 (4000-4999)
	ThirdPartyError ErrorCode = 4000 // 第三方服务通用错误
)

type Error struct {
	Code    ErrorCode
	Message string
	Data    interface{}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code, e.Message)
}

// WithData 添加错误详情
func (e *Error) WithData(data interface{}) *Error {
	e.Data = data
	return e
}

func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// IsError 判断错误是否为特定错误码
func IsError(err error, code ErrorCode) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*Error); ok {
		return e.Code == code
	}

	return false
}

// HTTPStatusFromError 根据错误获取 HTTP 状态码
func HTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	e, ok := err.(*Error)
	if !ok {
		return http.StatusInternalServerError
	}

	switch e.Code {
	case Success:
		return http.StatusOK
	case InvalidParameter:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case NotFound:
		return http.StatusNotFound
	case Timeout:
		return http.StatusRequestTimeout
	default:
		switch {
		case e.Code >= 1000 && e.Code < 2000:
			return http.StatusBadRequest // 业务错误
		case e.Code >= 2000 && e.Code < 4000:
			return http.StatusInternalServerError // 数据库、缓存错误
		case e.Code >= 4000 && e.Code < 5000:
			return http.StatusInternalServerError // 第三方服务错误
		default:
			return http.StatusInternalServerError // 其他未定义错误
		}
	}
}

// ErrorResponse 定义统一的错误响应结构
type ErrorResponse struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewErrorResponse 从错误创建错误响应
func NewErrorResponse(err error) *ErrorResponse {
	if err == nil {
		return &ErrorResponse{
			Code:    Success,
			Message: "success",
		}
	}

	if e, ok := err.(*Error); ok {
		return &ErrorResponse{
			Code:    e.Code,
			Message: e.Message,
			Data:    e.Data,
		}
	}

	return &ErrorResponse{
		Code:    UnknownError,
		Message: err.Error(),
	}
}
