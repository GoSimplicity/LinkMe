package domain

// User 表示用户信息的结构体
type User struct {
	ID          int64   // 用户ID，主键
	Phone       *string // 手机号码，指针类型，允许为空
	Email       string  // 邮箱地址，唯一
	Password    string  // 密码
	CreateTime  int64   // 创建时间，Unix时间戳
	UpdatedTime int64
	Deleted     bool    // 删除标志，表示该用户是否被删除
	Profile     Profile // 用户的详细资料
}

// Profile 表示用户详细资料的结构体
type Profile struct {
	ID       int64  // 资料ID，主键
	UserID   int64  // 用户ID，外键，关联到用户
	NickName string // 昵称
	Avatar   string // 头像URL
	About    string // 个人简介
	Birthday string // 生日
}

type UserWithProfileAndRule struct {
	ID           int64
	PasswordHash string
	Deleted      bool
	Email        string
	Phone        *string
	ProfileID    int64
	UserID       int64
	NickName     string
	Avatar       string
	About        string
	Birthday     string
	Role         string
}
