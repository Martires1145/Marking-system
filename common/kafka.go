package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"time"
)

var (
	reader *kafka.Reader
	topic1 = "test1"
	host   = viper.GetString("kafka.host")
)

func WriteTopicID(ctx context.Context, topicR, topicW string) {
	fmt.Println(topic1, topicR, topicW)
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(host),
		Topic:                  topic1,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireNone,
		Async:                  true,
		AllowAutoTopicCreation: true,
		BatchSize:              1,
	}
	r, _ := json.Marshal(topicR)
	w, _ := json.Marshal(topicW)

	defer writer.Close()
	for i := 0; i < 3; i++ {
		if err := writer.WriteMessages(
			ctx,
			kafka.Message{Key: []byte("1145"), Value: w},
		); err != nil {
			if err == kafka.LeaderNotAvailable {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				fmt.Printf("写入失败:%v\n", err)
			}
		} else {
			fmt.Println("topicW写入完成")
			break
		}
	}

	if err := writer.WriteMessages(
		ctx,
		kafka.Message{Key: []byte("1145"), Value: r},
	); err != nil {
		fmt.Printf("写入失败:%v\n", err)
	} else {
		fmt.Println("topicR写入完成")
	}
}

func ReadTopicAck(ctx context.Context, topic string) bool {
	fmt.Println(topic)
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{host},
		Topic:          topic,
		CommitInterval: 500 * time.Millisecond,
		StartOffset:    kafka.FirstOffset,
		Partition:      0,
	})
	defer reader.Close()

	if message, err := reader.ReadMessage(ctx); err != nil {
		fmt.Printf("读kafka失败:%v\n", err)
		return false
	} else {
		fmt.Println("收到ACK")
		if string(message.Value) == "ACK" {
			return true
		}
	}
	return false
}

func WriteTopicAck(ctx context.Context, topic string) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(host),
		Topic:                  topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireNone,
		Async:                  true,
		AllowAutoTopicCreation: true,
		BatchSize:              1,
	}
	defer writer.Close()
	for i := 0; i < 3; i++ {
		if err := writer.WriteMessages(
			ctx,
			kafka.Message{Key: []byte("1145"), Value: []byte("ACK")},
		); err != nil {
			if err == kafka.LeaderNotAvailable {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				fmt.Printf("写入失败:%v\n", err)
			}
		} else {
			fmt.Println("ACK写入完成")
			break
		}
	}
}

func WriteMsg(ctx context.Context, topic string, msg []byte) error {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(host),
		Topic:                  topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireNone,
		Async:                  true,
		AllowAutoTopicCreation: true,
		BatchSize:              1,
	}
	defer writer.Close()

	for i := 0; i < 3; i++ {
		if err := writer.WriteMessages(
			ctx,
			kafka.Message{Key: []byte("1145"), Value: msg},
		); err != nil {
			if err == kafka.LeaderNotAvailable {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				fmt.Printf("写入失败:%v\n", err)
				return err
			}
		} else {
			fmt.Println("msg写入完成")
			break
		}
	}
	return nil
}

func ReadMsg(ctx context.Context, topic string, channel *chan []byte) {
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:          []string{host},
		Topic:            topic,
		CommitInterval:   500 * time.Millisecond,
		StartOffset:      kafka.LastOffset,
		ReadBatchTimeout: 10 * time.Millisecond,
		Partition:        0,
	})
	defer reader.Close()

	for {
		if message, err := reader.ReadMessage(ctx); err != nil {
			fmt.Printf("读kafka失败:%v\n", err)
			break
		} else {
			*channel <- message.Value
		}
	}
}

func CreateTopic(ctx context.Context, topicR, topicW string) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(host),
		Topic:                  topicR,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireNone,
		Async:                  true,
		AllowAutoTopicCreation: true,
		BatchSize:              1,
	}

	_ = writer.WriteMessages(ctx, kafka.Message{Key: []byte("1145"), Value: []byte("create")})
	fmt.Println("topicR created")

	writer = &kafka.Writer{
		Addr:                   kafka.TCP(host),
		Topic:                  topicW,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireNone,
		Async:                  true,
		AllowAutoTopicCreation: true,
		BatchSize:              1,
	}

	_ = writer.WriteMessages(ctx, kafka.Message{Key: []byte("1145"), Value: []byte("create")})
	fmt.Println("topicW created")

}
