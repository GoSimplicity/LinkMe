package domain

import (
	"errors"
	"time"

	"github.com/dlclark/regexp2"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUsernameFormat = errors.New("invalid username format")
	ErrInvalidPasswordFormat = errors.New("invalid password format")
	ErrPasswordMismatch      = errors.New("password mismatch")
)

const (
	usernameRegexPattern = `^[a-zA-Z0-9]{6,}$` // 至少6位的字母数字组合
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type User struct {
	ID          int64   `json:"id"`          // 用户ID，主键
	Username    string  `json:"username"`    // 用户名，唯一
	Phone       *string `json:"phone"`       // 手机号码，指针类型，允许为空
	Email       string  `json:"email"`       // 邮箱地址，唯一
	Password    string  `json:"password"`    // 密码
	CreateTime  int64   `json:"createTime"`  // 创建时间，Unix时间戳
	UpdatedTime int64   `json:"updatedTime"` // 更新时间，Unix时间戳
	Deleted     bool    `json:"deleted"`     // 删除标志，表示该用户是否被删除
	Profile     Profile `json:"profile"`     // 用户的详细资料
}

type Profile struct {
	ID       int64  `json:"id"`       // 资料ID，主键
	UserID   int64  `json:"userId"`   // 用户ID，外键，关联到用户
	RealName string `json:"realName"` // 真实姓名
	Avatar   string `json:"avatar"`   // 头像URL
	Email    string `json:"email"`    // 邮箱
	About    string `json:"about"`    // 个人简介
	Birthday string `json:"birthday"` // 生日
}

type UserWithProfile struct {
	ID           int64   `json:"id"`
	Username     string  `json:"username"`
	PasswordHash string  `json:"passwordHash"`
	Deleted      bool    `json:"deleted"`
	Phone        *string `json:"phone"`
	Email        string  `json:"email"`
	ProfileID    int64   `json:"profileId"`
	UserID       int64   `json:"userId"`
	RealName     string  `json:"realName"`
	Avatar       string  `json:"avatar"`
	About        string  `json:"about"`
	Birthday     string  `json:"birthday"`
}

// ValidateUsername 验证用户名格式
func (u *User) ValidateUsername() error {
	usernameRegex := regexp2.MustCompile(usernameRegexPattern, regexp2.None)
	isMatch, err := usernameRegex.MatchString(u.Username)
	if err != nil {
		return err
	}
	if !isMatch {
		return ErrInvalidUsernameFormat
	}
	return nil
}

// ValidatePassword 验证密码格式
func (u *User) ValidatePassword() error {
	passwordRegex := regexp2.MustCompile(passwordRegexPattern, regexp2.None)
	isMatch, err := passwordRegex.MatchString(u.Password)
	if err != nil {
		return err
	}
	if !isMatch {
		return ErrInvalidPasswordFormat
	}
	return nil
}

// VerifyPassword 验证密码是否匹配
func (u *User) VerifyPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return ErrPasswordMismatch
	}
	return nil
}

// HashPassword 对密码进行哈希处理
func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// UpdateProfile 更新用户资料
func (u *User) UpdateProfile(newProfile Profile) {
	u.Profile = newProfile
	u.UpdatedTime = time.Now().Unix()
}

// MarkAsDeleted 标记用户为已删除
func (u *User) MarkAsDeleted() {
	u.Deleted = true
	u.UpdatedTime = time.Now().Unix()
}
