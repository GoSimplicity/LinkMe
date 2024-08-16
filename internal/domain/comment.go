package domain

type Comment struct {
	Id            int64
	UserId        int64
	Biz           string
	BizId         int64
	PostId        int64
	Content       string
	RootComment   *Comment  // 根节点
	ParentComment *Comment  // 父节点
	Children      []Comment // 子节点
	CreatedAt     int64
	UpdatedAt     int64
}
