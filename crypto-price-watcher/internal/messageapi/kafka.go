package messageapi

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
)

var (
	//TODO: get kafka host from config
	kafkaHost string = "localhost:9092"
	logger           = watermill.NewStdLogger(false, true)
)

func NewPubliser() (*kafka.Publisher, error) {
	config := kafka.PublisherConfig{
		Brokers:   []string{kafkaHost},
		Marshaler: kafka.DefaultMarshaler{},
	}
	return kafka.NewPublisher(config, logger)
}

func NewSubscriber() (*kafka.Subscriber, error) {
	config := kafka.SubscriberConfig{
		Brokers:               []string{kafkaHost},
		Unmarshaler:           kafka.DefaultMarshaler{},
		OverwriteSaramaConfig: kafka.DefaultSaramaSubscriberConfig(),
	}
	return kafka.NewSubscriber(config, logger)
}
