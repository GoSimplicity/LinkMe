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

package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// SHA256 计算字符串的SHA-256哈希值
// 注意：这是基础哈希函数，不应直接用于密码存储，推荐使用GeneratePasswordHash
func SHA256(data string) string {
	if data == "" {
		return ""
	}
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateSalt 生成随机盐值
// size: 盐值的字节长度
func GenerateSalt(size int) (string, error) {
	if size <= 0 {
		return "", errors.New("salt size must be positive")
	}

	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to generate random salt: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}

// SHA256WithSalt 使用给定盐值计算字符串的SHA-256哈希值
// 注意：这是基础哈希+盐函数，不应直接用于密码存储，推荐使用GeneratePasswordHash
func SHA256WithSalt(data string, salt string) string {
	if data == "" || salt == "" {
		return ""
	}

	hash := sha256.New()
	hash.Write([]byte(data + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

// GeneratePasswordHash 生成安全的密码哈希
// 返回格式: "$sha256$salt$hash"，其中salt是随机生成的盐值
func GeneratePasswordHash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// 生成32字节(256位)的随机盐值
	salt, err := GenerateSalt(32)
	if err != nil {
		return "", err
	}

	// 计算哈希
	hash := SHA256WithSalt(password, salt)

	// 返回格式化的哈希字符串
	return fmt.Sprintf("$sha256$%s$%s", salt, hash), nil
}

// VerifyPasswordHash 验证密码是否匹配存储的哈希值
// storedHash: 之前由GeneratePasswordHash生成的哈希字符串
// password: 待验证的密码
func VerifyPasswordHash(storedHash string, password string) (bool, error) {
	if storedHash == "" || password == "" {
		return false, errors.New("hash and password cannot be empty")
	}

	// 解析存储的哈希字符串
	parts := strings.Split(storedHash, "$")
	if len(parts) != 4 || parts[0] != "" || parts[1] != "sha256" {
		return false, errors.New("invalid hash format")
	}

	salt := parts[2]
	expectedHash := parts[3]

	// 计算提供的密码的哈希值
	computedHash := SHA256WithSalt(password, salt)

	// 使用时间恒定比较函数防止时序攻击
	return subtle.ConstantTimeCompare([]byte(computedHash), []byte(expectedHash)) == 1, nil
}

// IsValidHash 检查哈希字符串格式是否有效
func IsValidHash(hash string) bool {
	parts := strings.Split(hash, "$")
	return len(parts) == 4 && parts[0] == "" && parts[1] == "sha256" && parts[2] != "" && parts[3] != ""
}

// GenerateSecureToken 生成安全随机令牌(可用于会话ID、CSRF令牌等)
func GenerateSecureToken(byteLength int) (string, error) {
	if byteLength <= 0 {
		return "", errors.New("byte length must be positive")
	}

	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
