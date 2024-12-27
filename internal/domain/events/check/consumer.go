package check

import (
	"context"
	"encoding/json"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type CheckEventConsumer struct {
	checkRepo repository.CheckRepository
	client    sarama.Client
	l         *zap.Logger
}

func NewCheckEventConsumer(
	checkRepo repository.CheckRepository,
	client sarama.Client,
	l *zap.Logger,
) *CheckEventConsumer {
	return &CheckEventConsumer{
		checkRepo: checkRepo,
		client:    client,
		l:         l,
	}
}

func (i *CheckEventConsumer) Start(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroupFromClient("check", i.client)
	if err != nil {
		i.l.Error("创建消费者组失败", zap.Error(err))
		return err
	}

	i.l.Info("CheckConsumer 开始消费")

	go func() {
		defer cg.Close()
		for {
			select {
			case <-ctx.Done():
				i.l.Info("审核消息消费者停止")
				return
			default:
				if err := cg.Consume(ctx, []string{TopicCheckEvent}, &checkConsumerGroupHandler{r: i}); err != nil {
					i.l.Error("审核消息消费循环出错", zap.Error(err))
					continue
				}
			}
		}
	}()

	return nil
}

type checkConsumerGroupHandler struct {
	r *CheckEventConsumer
}

func (h *checkConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *checkConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *checkConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.r.ConsumeCheck(sess, msg)
	}
	return nil
}

// ConsumeCheck 处理单条审核消息
func (i *CheckEventConsumer) ConsumeCheck(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	var evt CheckEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return
	}

	// 参数校验
	if evt.PostId == 0 || evt.Uid == 0 {
		i.l.Warn("无效的审核事件",
			zap.Uint("post_id", evt.PostId),
			zap.Int64("uid", evt.Uid))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	check := domain.Check{
		PostID:  evt.PostId,
		Uid:     evt.Uid,
		Title:   evt.Title,
		Content: evt.Content,
	}

	// 创建审核记录
	if _, err := i.checkRepo.Create(ctx, check); err != nil {
		i.l.Error("创建审核记录失败", zap.Error(err))
		return
	}

	i.l.Info("创建审核记录成功",
		zap.Uint("post_id", evt.PostId),
		zap.Int64("uid", evt.Uid))

	sess.MarkMessage(msg, "")
}
