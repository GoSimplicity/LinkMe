package AiCheck

import (
	"context"
	"fmt"
	ark "github.com/sashabaranov/go-openai"
	"testing"
)

func TestContent(t *testing.T) {
	ARK_API_KEY := "b3977816-2a07-44df-9fe2-4ec02224e147"
	config := ark.DefaultConfig(ARK_API_KEY)
	config.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	client := ark.NewClientWithConfig(config)

	//fmt.Println("----- standard request -----")
	checkLanguage := "你是一个负责评论审核的人工智能，请判断输入内容是否包含违规信息,返回结果1或者0，其中1表示包含，0表示不包含。"
	//checkContents := "操你妈hhh"
	checkContents := "hhh"
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

	}
	fmt.Println(resp.Choices[0].Message.Content)
}
