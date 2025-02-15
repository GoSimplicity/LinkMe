package canal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

// Debezium 变更事件结构
type ChangeEvent struct {
	Schema struct {
		Type   string `json:"type"`
		Fields []struct {
			Type   string `json:"type"`
			Fields []struct {
				Type       string      `json:"type"`
				Optional   bool        `json:"optional"`
				Field      string      `json:"field"`
				Default    interface{} `json:"default,omitempty"`
				Name       string      `json:"name,omitempty"`
				Version    int         `json:"version,omitempty"`
				Parameters struct {
					Allowed string `json:"allowed"`
				} `json:"parameters,omitempty"`
			} `json:"fields,omitempty"`
			Optional bool   `json:"optional"`
			Name     string `json:"name,omitempty"`
			Field    string `json:"field"`
			Version  int    `json:"version,omitempty"`
		} `json:"fields"`
		Optional bool   `json:"optional"`
		Name     string `json:"name"`
		Version  int    `json:"version"`
	} `json:"schema"`
	Payload struct { //数据的主体部分
		Before struct { //数据操作前
			Id           int         `json:"id"`
			CreatedAt    interface{} `json:"created_at"`
			UpdatedAt    interface{} `json:"updated_at"`
			DeletedAt    interface{} `json:"deleted_at"`
			PasswordHash string      `json:"password_hash"`
			Deleted      int         `json:"deleted"`
			Email        string      `json:"email"`
			Phone        string      `json:"phone"`
			Nickname     string      `json:"nickname"`
		} `json:"before"`
		After struct { //数据操作后
			Id           int         `json:"id"`
			CreatedAt    interface{} `json:"created_at"`
			UpdatedAt    interface{} `json:"updated_at"`
			DeletedAt    interface{} `json:"deleted_at"`
			PasswordHash string      `json:"password_hash"`
			Deleted      int         `json:"deleted"`
			Email        string      `json:"email"`
			Phone        string      `json:"phone"`
			Nickname     string      `json:"nickname"`
		} `json:"after"`
		Source struct {
			Version   string      `json:"version"`
			Connector string      `json:"connector"` //创建链接器的对象即被备份的数据库
			Name      string      `json:"name"`      //自定义topic前缀
			TsMs      int64       `json:"ts_ms"`
			Snapshot  string      `json:"snapshot"` //是否是快照
			Db        string      `json:"db"`       //数据库名
			Sequence  interface{} `json:"sequence"`
			TsUs      int64       `json:"ts_us"`
			TsNs      int64       `json:"ts_ns"`
			Table     string      `json:"table"` //表名
			ServerId  int         `json:"server_id"`
			Gtid      interface{} `json:"gtid"`
			File      string      `json:"file"` //从哪个binlog日志中取得数据的
			Pos       int         `json:"pos"`  //偏移量
			Row       int         `json:"row"`
			Thread    int         `json:"thread"`
			Query     interface{} `json:"query"`
		} `json:"source"`
		Transaction interface{} `json:"transaction"`
		Op          string      `json:"op"` //指令类型(c:INSERT u:UPDATE d:DELETE
		TsMs        int64       `json:"ts_ms"`
		TsUs        int64       `json:"ts_us"`
		TsNs        int64       `json:"ts_ns"`
	} `json:"payload"`
}

// consumerGroupHandler 结构体实现了Kafka Consumer Group的接口
type consumerGroupHandler struct {
}

func initKafka() sarama.Client {
	scfg := sarama.NewConfig()
	scfg.Consumer.Offsets.Initial = sarama.OffsetOldest // 关键配置：从头消费
	client, _ := sarama.NewClient([]string{"192.168.84.130:9092"}, scfg)
	return client
}

func TestConsumer(t *testing.T) {
	// 配置 Kafka 消费者
	cg, err := sarama.NewConsumerGroupFromClient("es_consumer_group", initKafka())
	if err != nil {
		panic(err)
	}

	// 订阅 Topic
	go func() {
		for {
			if err := cg.Consume(context.Background(), []string{"oracle.linkme.users"}, &consumerGroupHandler{}); err != nil {
				fmt.Printf("退出了消费循环异常: %v", err)
				time.Sleep(time.Second * 5)
			}
			time.Sleep(time.Second * 2)
		}
	}()

	// 优雅退出
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigchan:
		fmt.Printf("consumer test terminated")
	}

}

func handleEvent(event ChangeEvent) {
	switch event.Payload.Op {
	case "c":
		fmt.Printf("INSERT: %+v\n", event.Payload.After)
	case "u":
		fmt.Printf("UPDATE: Before=%+v, After=%+v\n", event.Payload.Before, event.Payload.After)
	case "d":
		fmt.Printf("DELETE: %+v\n", event.Payload.Before)
	default:
		fmt.Printf("Unknown operation: %s\n", event.Payload.Op)
	}
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim 消费Kafka的消息
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// 消费消息
		for {

			var event ChangeEvent

			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Printf("Failed to parse event: %v\n", err)
				continue
			}

			//fmt.Printf("event: %v\n\n", event)

			// 处理事件
			handleEvent(event)

		}
	}
	return nil
}
