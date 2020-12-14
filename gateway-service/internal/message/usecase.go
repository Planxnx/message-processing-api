package message

import (
	"fmt"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/protobuf/proto"
)

type MessageUsecase struct {
	KafkaPublisher *kafka.Publisher
}

func New(k *kafka.Publisher) *MessageUsecase {
	return &MessageUsecase{
		KafkaPublisher: k,
	}
}

func (m *MessageUsecase) EmitCommon(uuid string, msg *messageschema.DefaultMessage) error {
	msgByte, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	kafkaMsg := message.NewMessage(uuid, msgByte)
	if err := m.KafkaPublisher.Publish(messageschema.CommonMessageTopic, kafkaMsg); err != nil {
		return fmt.Errorf("failed on publish message: %v", err)
	}

	return nil
}

func (m *MessageUsecase) EmitReply(uuid string, msg *messageschema.DefaultMessage) error {
	msgByte, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	kafkaMsg := message.NewMessage(uuid, msgByte)
	if err := m.KafkaPublisher.Publish(messageschema.ReplyMessageTopic, kafkaMsg); err != nil {
		return fmt.Errorf("failed on publish message: %v", err)
	}

	return nil
}

func (m *MessageUsecase) Emit(uuid string, topic string, msg *messageschema.DefaultMessage) error {
	msgByte, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	kafkaMsg := message.NewMessage(uuid, msgByte)
	if err := m.KafkaPublisher.Publish(topic, kafkaMsg); err != nil {
		return fmt.Errorf("failed on publish message topic %s : %v", topic, err)
	}

	return nil
}
