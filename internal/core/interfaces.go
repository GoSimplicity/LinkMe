package core

// 通用 Repository 接口
type Repository interface {
	Create(entity interface{}) error
	FindByID(id interface{}) (interface{}, error)
	FindOne(query interface{}) (interface{}, error)
	FindAll(conditions map[string]interface{}, page, size int) ([]interface{}, int64, error)
	Update(entity interface{}) error
	Delete(id interface{}) error
	Transaction(fn func() error) error
}

// 通用 Service 接口
type Service interface {
	Create(ctx interface{}, dto interface{}) (interface{}, error)
	Get(ctx interface{}, query interface{}) (interface{}, error)
	List(ctx interface{}, query interface{}, page, size int) ([]interface{}, int64, error)
	Update(ctx interface{}, dto interface{}) error
	Delete(ctx interface{}, id interface{}) error
}

// 通用 Database 接口
type Database interface {
	Create(ctx interface{}, entity interface{}) error
	FindOne(ctx interface{}, query interface{}) (interface{}, error)
	FindAll(ctx interface{}, query interface{}, page, size int) ([]interface{}, int64, error)
	Update(ctx interface{}, query interface{}, updates interface{}) error
	Delete(ctx interface{}, query interface{}) error
	Transaction(ctx interface{}, fn func(ctx interface{}) error) error
	ExecRaw(query string, args ...interface{}) (interface{}, error)
	QueryRaw(query string, args ...interface{}) ([]interface{}, error)
}
