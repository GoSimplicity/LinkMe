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

package migrations

import (
	"database/sql"
	"errors"
	"time"

	"github.com/GoSimplicity/LinkMe/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

// User 用户模型
type User struct {
	ID        uint64                `gorm:"primaryKey;autoIncrement"`                             // 用户ID
	Username  string                `gorm:"type:varchar(50);uniqueIndex:udx_deleted_at;not null"` // 用户名(唯一)
	Password  string                `gorm:"type:varchar(255);not null"`                           // 密码(加密存储)
	RealName  string                `gorm:"type:varchar(50)"`                                     // 真实姓名
	Avatar    string                `gorm:"type:varchar(255)"`                                    // 头像URL
	About     string                `gorm:"type:text"`                                            // 个人简介
	Birthday  sql.NullTime          `gorm:"type:date"`                                            // 生日
	Email     string                `gorm:"type:varchar(100);uniqueIndex:udx_deleted_at"`         // 邮箱(唯一)
	Phone     string                `gorm:"type:varchar(20);uniqueIndex:udx_deleted_at"`          // 手机号(唯一)
	Status    int8                  `gorm:"type:tinyint;default:1;not null"`                      // 用户状态(0:禁用,1:启用)
	Role      string                `gorm:"type:varchar(20);default:'user';not null"`             // 用户角色(user/admin)
	LastLogin sql.NullTime          `gorm:"type:datetime"`                                        // 最后登录时间
	CreatedAt time.Time             `gorm:"not null"`                                             // 创建时间
	UpdatedAt time.Time             `gorm:"not null"`                                             // 更新时间
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:udx_deleted_at"`                           // 删除时间(软删除)
}

// TableName 指定表名
func (User) TableName() string {
	return "lm_users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	if !utils.IsValidHash(u.Password) {
		// 对密码进行加密
		hashPassword, err := utils.GeneratePasswordHash(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashPassword
	}
	return
}

// BeforeUpdate 更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	if u.Password != "" {
		if !utils.IsValidHash(u.Password) {
			hashPassword, err := utils.GeneratePasswordHash(u.Password)
			if err != nil {
				return err
			}
			u.Password = hashPassword
		}
	}
	return
}

// BeforeDelete 删除前钩子
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.Role == "admin" {
		return errors.New("admin user not allowed to delete")
	}
	return
}
