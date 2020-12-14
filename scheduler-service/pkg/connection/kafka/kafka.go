package kafka

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
)

const kafkaHost string = "localhost:9092"

func NewPubliser() (*kafka.Publisher, error) {
	config := kafka.PublisherConfig{
		Brokers:   []string{kafkaHost},
		Marshaler: kafka.DefaultMarshaler{},
	}
	return kafka.NewPublisher(config, watermill.NewStdLogger(false, true))
}
