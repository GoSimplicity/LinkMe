package domain

type IM struct {
	Conversations []Conversation // 用户的会话列表
}

// Message 表示会话中的单条消息
type Message struct {
	ID             int64  // 消息ID，主键
	ConversationID int64  // 会话ID，外键，关联到对应的会话
	SenderID       int64  // 发送消息的用户ID，外键
	Content        string // 消息内容
	Timestamp      int64  // 消息发送时间，Unix时间戳
	IsRead         bool   // 消息的已读状态
}

// Conversation 用户之间的聊天会话
type Conversation struct {
	ID           int64         // 会话ID，主键
	Participants []Participant // 参与会话的用户列表
	Messages     []Message     // 会话中的消息列表
	LastMessage  Message       // 会话中的最后一条消息
	CreateTime   int64         // 会话创建时间，Unix时间戳
	UpdatedTime  int64         // 会话最后更新时间，Unix时间戳
}

// Participant 参与会话的用户
type Participant struct {
	ID         int64 // 参与者ID，主键
	UserID     int64 // 用户ID，外键，关联到用户
	IsMuted    bool  // 用户是否静音此会话
	IsDeleted  bool  // 用户是否删除了此会话
	JoinedTime int64 // 用户加入会话的时间，Unix时间戳
}
