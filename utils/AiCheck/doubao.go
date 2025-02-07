package AiCheck

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	ark "github.com/sashabaranov/go-openai"
	"github.com/sony/gobreaker"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	client     *ark.Client
	clientOnce sync.Once
	breaker    *gobreaker.CircuitBreaker
	logger     *zap.Logger
)

func init() {
	breaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "AI_Check",
		Timeout:     3 * time.Second,    // 超时时间
		MaxRequests: 5,                  // 半开状态下允许的试探请求
		Interval:    1 * time.Minute,    // 统计窗口间隔
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5 // 连续失败触发熔断
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			logger.Info("熔断状态变更", zap.String("name", name), zap.Any("from", from), zap.Any("to", to))
		},
	})
}

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

func CheckPostContent(title string, content string) (bool, error) {
	result, err := breaker.Execute(func() (interface{}, error) {
		client := getClient()
		checkLanguage := "你是一个负责评论审核的人工智能，请对输入的内容进行审查，判断其是否包含违规信息。返回结果如下: 如果是 1 说明内容包含违规信息, 如果是 0 说明内容不包含违规信息, 如果是 -1 说明无法判断或其他错误"
		checkContents := fmt.Sprintf("标题:%s\n内容:%s", title, content)
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
		ans := true
		if resp.Choices[0].Message.Content == "1" || resp.Choices[0].Message.Content == "-1" {
			ans = false
		}
		return ans, nil
	})

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}
