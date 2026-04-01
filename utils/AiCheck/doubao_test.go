//go:build integration
// +build integration

package AiCheck

import (
	"context"
	"fmt"
	ark "github.com/sashabaranov/go-openai"
	"os"
	"testing"
)

func TestContent(t *testing.T) {
	if os.Getenv("LINKME_ARK_API_PROVIDER") != "ark" {
		t.Skip("当前未启用 Ark 审核 provider")
	}
	ARK_API_KEY := os.Getenv("LINKME_ARK_API_KEY")
	if ARK_API_KEY == "" {
		t.Skip("Ark API Key 未配置")
	}
	cfg := ark.DefaultConfig(ARK_API_KEY)
	cfg.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	client := ark.NewClientWithConfig(cfg)

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
		t.Fatalf("ChatCompletion error: %v", err)
	}
	if len(resp.Choices) == 0 {
		t.Fatal("Ark 返回结果为空")
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
