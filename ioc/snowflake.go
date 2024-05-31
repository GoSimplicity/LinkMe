package ioc

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

// InitializeSnowflakeNode 初始化雪花节点
func InitializeSnowflakeNode() *sf.Node {
	st := time.Now()                   // 直接使用当前时间
	sf.Epoch = st.UnixNano() / 1000000 // 计算起始时间戳并赋值给sf.Epoch
	node, err := sf.NewNode(1)         // 创建新的节点
	if err != nil {
		return nil
	}
	return node
}
