package canal

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"os"
	"os/signal"
	"syscall"
	"testing"
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
	Payload struct {
		Before struct {
			Id           int         `json:"id"`
			CreatedAt    interface{} `json:"created_at"`
			UpdatedAt    interface{} `json:"updated_at"`
			DeletedAt    interface{} `json:"deleted_at"`
			PasswordHash string      `json:"password_hash"`
			Deleted      int         `json:"deleted"`
			Email        string      `json:"email"`
			Phone        string      `json:"phone"`
		} `json:"before"`
		After struct {
			Id           int         `json:"id"`
			CreatedAt    interface{} `json:"created_at"`
			UpdatedAt    interface{} `json:"updated_at"`
			DeletedAt    interface{} `json:"deleted_at"`
			PasswordHash string      `json:"password_hash"`
			Deleted      int         `json:"deleted"`
			Email        string      `json:"email"`
			Phone        string      `json:"phone"`
		} `json:"after"`
		Source struct {
			Version   string      `json:"version"`
			Connector string      `json:"connector"`
			Name      string      `json:"name"`
			TsMs      int64       `json:"ts_ms"`
			Snapshot  string      `json:"snapshot"`
			Db        string      `json:"db"`
			Sequence  interface{} `json:"sequence"`
			TsUs      int64       `json:"ts_us"`
			TsNs      int64       `json:"ts_ns"`
			Table     string      `json:"table"`
			ServerId  int         `json:"server_id"`
			Gtid      interface{} `json:"gtid"`
			File      string      `json:"file"`
			Pos       int         `json:"pos"`
			Row       int         `json:"row"`
			Thread    int         `json:"thread"`
			Query     interface{} `json:"query"`
		} `json:"source"`
		Transaction interface{} `json:"transaction"`
		Op          string      `json:"op"`
		TsMs        int64       `json:"ts_ms"`
		TsUs        int64       `json:"ts_us"`
		TsNs        int64       `json:"ts_ns"`
	} `json:"payload"`
}

func TestConsumer(t *testing.T) {
	// 配置 Kafka 消费者
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// 创建消费者
	consumer, err := sarama.NewConsumer([]string{"192.168.84.130:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// 订阅 Topic
	partitionConsumer, err := consumer.ConsumePartition("oracle.linkme.users", 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	defer partitionConsumer.Close()

	// 优雅退出
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// 消费消息
	for {
		select {
		case <-sigchan:
			fmt.Println("Exiting...")
			return
		case msg := <-partitionConsumer.Messages():

			var event ChangeEvent

			fmt.Println("value:", string(msg.Value))
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Printf("Failed to parse event: %v\n", err)
				continue
			}

			// 处理事件
			handleEvent(event)
		}
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
