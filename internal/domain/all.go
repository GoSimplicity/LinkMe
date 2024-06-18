package domain

type Pagination struct {
	Page int    // 当前页码
	Size *int64 // 每页数据
	Uid  int64
	// 以下字段通常在服务端内部使用，不需要客户端传递
	Offset *int64 // 数据偏移量
	Total  *int64 // 总数据量
}
