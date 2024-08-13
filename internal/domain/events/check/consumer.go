package check

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/IBM/sarama"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"reflect"
	"time"
)

type CheckConsumer struct {
	client sarama.Client
	l      *zap.Logger
	repo   repository.PostRepository
}

type Event struct {
	Type     string                   `json:"type"`
	Database string                   `json:"database"`
	Table    string                   `json:"table"`
	Data     []map[string]interface{} `json:"data"`
}

type Check struct {
	ID        int64  `mapstructure:"id"`
	PostID    uint   `mapstructure:"post_id"`
	Title     string `mapstructure:"title"`
	Content   string `mapstructure:"content"`
	CreatedAt int64  `mapstructure:"created_at"`
	UpdatedAt int64  `mapstructure:"updated_at"`
	UserID    int64  `mapstructure:"author_id"`
	Status    uint8  `mapstructure:"status"`
	Remark    string `mapstructure:"remark"`
}

type consumerGroupHandler struct {
	r *CheckConsumer
}

func NewSyncConsumer(client sarama.Client, l *zap.Logger, repo repository.PostRepository) *CheckConsumer {
	return &CheckConsumer{
		client: client,
		l:      l,
		repo:   repo,
	}
}

func (r *CheckConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("check_consumer_group", r.client)
	if err != nil {
		return err
	}

	go func() {
		for {
			if err := cg.Consume(ctx, []string{"linkme_binlog"}, &consumerGroupHandler{r: r}); err != nil {
				r.l.Error("退出了消费循环异常", zap.Error(err))
				time.Sleep(time.Second * 5)
			}
		}
	}()

	return nil
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.r.Consume(sess, msg); err != nil {
			h.r.l.Error("处理消息失败", zap.Error(err), zap.ByteString("message", msg.Value))
		} else {
			sess.MarkMessage(msg, "")
		}
	}

	return nil
}

func (r *CheckConsumer) Consume(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) error {
	var e Event
	var checks []Check

	// 反序列化消息
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		return err
	}

	// 确保处理的表是目标表
	if e.Table != "checks" {
		return nil
	}

	// 将事件数据映射到 Check 结构体切片
	if err := decodeEventDataToChecks(e.Data, &checks); err != nil {
		return err
	}

	// 处理每个 Check 实例
	for _, check := range checks {
		if err := r.handlePost(sess.Context(), check); err != nil {
			return err
		}
	}

	return nil
}

// handlePost 根据状态处理帖子
func (r *CheckConsumer) handlePost(ctx context.Context, check Check) error {
	if check.Status == domain.Approved {
		return r.repo.UpdateStatus(ctx, domain.Post{
			ID:     check.PostID,
			Status: domain.Published,
		})
	}

	return nil
}

// decodeEventDataToChecks 将事件数据解码到 Check 结构体切片
func decodeEventDataToChecks(data []map[string]interface{}, checks *[]Check) error {
	config := &mapstructure.DecoderConfig{
		Result:           checks,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToTimeHookFunc("2006-01-02 15:04:05"),
			stringToNullTimeHookFunc("2006-01-02 15:04:05"),
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(data)
}

// stringToTimeHookFunc 将字符串转换为时间类型
func stringToTimeHookFunc(layout string) mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
		if f != reflect.String || t != reflect.Struct {
			return data, nil
		}

		str, ok := data.(string)
		if !ok || str == "" {
			return time.Time{}, nil
		}

		return time.Parse(layout, str)
	}
}

// stringToNullTimeHookFunc 将字符串转换为 sql.NullTime 类型
func stringToNullTimeHookFunc(layout string) mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
		if f != reflect.String || t != reflect.Struct {
			return data, nil
		}

		str, ok := data.(string)
		if !ok || str == "" {
			return sql.NullTime{Valid: false}, nil
		}

		parsedTime, err := time.Parse(layout, str)
		if err != nil {
			return nil, err
		}
		return sql.NullTime{Time: parsedTime, Valid: true}, nil
	}
}
