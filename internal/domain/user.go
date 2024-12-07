package domain

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmailFormat    = errors.New("invalid email format")
	ErrInvalidPasswordFormat = errors.New("invalid password format")
	ErrPasswordMismatch      = errors.New("password mismatch")
)

const (
	emailRegexPattern    = `^[a-zA-Z0-9]{6,}$` // 至少6位的字母数字组合
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type User struct {
	ID          int64   // 用户ID，主键
	Username    string  // 用户名，唯一
	Phone       *string // 手机号码，指针类型，允许为空
	Email       string  // 邮箱地址，唯一
	Password    string  // 密码
	CreateTime  int64   // 创建时间，Unix时间戳
	UpdatedTime int64
	Deleted     bool    // 删除标志，表示该用户是否被删除
	Profile     Profile // 用户的详细资料
}

// ValidateEmail 验证邮箱格式
func (u *User) ValidateEmail() error {
	emailRegex := regexp.MustCompile(emailRegexPattern)
	if !emailRegex.MatchString(u.Username) {
		return ErrInvalidEmailFormat
	}
	return nil
}

// ValidatePassword 验证密码格式
func (u *User) ValidatePassword() error {
	passwordRegex := regexp.MustCompile(passwordRegexPattern)
	if !passwordRegex.MatchString(u.Password) {
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

type Profile struct {
	ID       int64  // 资料ID，主键
	UserID   int64  // 用户ID，外键，关联到用户
	NickName string // 昵称
	Avatar   string // 头像URL
	Email    string // 邮箱
	About    string // 个人简介
	Birthday string // 生日
}

type UserWithProfileAndRule struct {
	ID           int64
	Username     string
	PasswordHash string
	Deleted      bool
	Phone        *string
	Email        string
	ProfileID    int64
	UserID       int64
	NickName     string
	Avatar       string
	About        string
	Birthday     string
	Role         string
}
