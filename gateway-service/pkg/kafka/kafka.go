package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
)

var (
	KafkaHost string = "localhost:9092"
	logger           = watermill.NewStdLogger(false, true)
)

func NewPubliser() (*kafka.Publisher, error) {
	kafka.DefaultSaramaSubscriberConfig()
	config := kafka.PublisherConfig{
		Brokers:   []string{KafkaHost},
		Marshaler: kafka.DefaultMarshaler{},
	}
	return kafka.NewPublisher(config, logger)
}

func NewSubscriber() (*kafka.Subscriber, error) {
	config := kafka.SubscriberConfig{
		Brokers:     []string{KafkaHost},
		Unmarshaler: kafka.DefaultMarshaler{},
		// ConsumerGroup:         "main",
	}
	return kafka.NewSubscriber(config, logger)
}

func NewSarama() (*sarama.Broker, error) {
	broker := sarama.NewBroker(KafkaHost)

	saramaConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	broker.Open(saramaConfig)

	if _, err := broker.Connected(); err != nil {
		return nil, err
	}

	return broker, nil
}

func CreateTopic(topic string, broker *sarama.Broker) error {
	topicDetails := map[string]*sarama.TopicDetail{
		topic: {
			NumPartitions:     1,
			ReplicationFactor: 1,
			ConfigEntries:     make(map[string]*string),
		},
	}

	_, err := broker.CreateTopics(&sarama.CreateTopicsRequest{
		Timeout:      time.Second * 15,
		TopicDetails: topicDetails,
	})

	return err
}

func DeleteTopic(topics []string, broker *sarama.Broker) error {
	_, err := broker.DeleteTopics(&sarama.DeleteTopicsRequest{
		Timeout: time.Second * 15,
		Topics:  topics,
	})

	return err
}
