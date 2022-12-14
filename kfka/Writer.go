package kfka

import (
	"context"
	"os"

	kafka "github.com/segmentio/kafka-go"
)

func Writer(topic string, key1 string, val1 string) {
	// get kafka writer using environment variables.
	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("KafkaURL")),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 1,
	}

	w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key1),
			Value: []byte(val1),
		},
	)

	w.Close()
}
