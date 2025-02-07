package AiCheck

import (
	"context"
	"errors"
	"fmt"
	ark "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"sync"
)

var (
	client     *ark.Client
	clientOnce sync.Once
)

type config struct {
	KEY string `yaml:"dsn"`
}

// 获取单例 client
func getClient() *ark.Client {
	clientOnce.Do(func() {
		var c config
		if err := viper.UnmarshalKey("ark_api", &c); err != nil {
			panic(fmt.Errorf("init failed：%v", err))
		}
		ARK_API_KEY := c.KEY // 这个AIP是可以修改的
		config := ark.DefaultConfig(ARK_API_KEY)
		config.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"
		client = ark.NewClientWithConfig(config)
		fmt.Println("AI 审核 client 初始化完成")
	})
	return client
}

func CheckPostContent(content string) (bool, error) {
	client := getClient()
	checkLanguage := "你是一个负责评论审核的人工智能，请判断输入内容是否包含违规信息,返回结果1或者0，其中1表示包含，0表示不包含。"
	//checkContents := "操你妈hhh"
	//checkContents := "hhh"
	checkContents := content
	resp, err := client.CreateChatCompletion(
		context.Background(),
		ark.ChatCompletionRequest{
			Model: "ep-20250207162731-kvrzk",
			Messages: []ark.ChatCompletionMessage{
				{
					Role:    ark.ChatMessageRoleSystem,
					Content: checkLanguage,
				},
				{
					Role:    ark.ChatMessageRoleUser,
					Content: checkContents,
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return false, errors.New("AI 审核失败")
	}
	//fmt.Println(resp.Choices[0].Message.Content)
	ans := true
	if resp.Choices[0].Message.Content == "1" {
		ans = false
	}
	return ans, nil
}
