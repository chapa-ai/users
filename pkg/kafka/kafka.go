package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
	users "users/grpc-users"
)

const (
	topic         = "topicNew191"
	brokerAddress = "localhost:9092"
)

var kaf *kafka.Conn

func InitKafka() (*kafka.Conn, error) {
	if kaf == nil {
		k, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, 0)
		if err != nil {
			log.Fatal("failed to dial leader:", err)
			return nil, err
		}
		kaf = k
	}
	return kaf, nil
}

func Kafka(ctx context.Context, insertedUser *users.User) (string, string, error) {
	kaf, err := InitKafka()
	if err != nil {
		return "", "", err
	}

	kaf.SetWriteDeadline(time.Now().Add(10 * time.Second))
	kaf.WriteMessages(kafka.Message{Key: []byte(insertedUser.Name), Value: []byte(insertedUser.Age)})

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})

	msg, err := r.ReadMessage(ctx)
	if err != nil {
		panic("could not read message " + err.Error())
	}

	key := string(msg.Key)
	val := string(msg.Value)

	return key, val, nil
}
