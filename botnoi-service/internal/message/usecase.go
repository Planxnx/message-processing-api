package message

import (
	"encoding/json"
	"fmt"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageUsecase struct {
	KafkaPublisher *kafka.Publisher
}

func NewUsecase(k *kafka.Publisher) *MessageUsecase {
	return &MessageUsecase{
		KafkaPublisher: k,
	}
}

func (m *MessageUsecase) EmitCommon(uuid string, msg *messageschema.DefaultMessageFormat) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	kafkaMsg := message.NewMessage(uuid, msgJSON)
	if err := m.KafkaPublisher.Publish(messageschema.CommonMessage, kafkaMsg); err != nil {
		return fmt.Errorf("failed on publish message: %v", err)
	}

	return nil
}

func (m *MessageUsecase) EmitReply(uuid string, msg *messageschema.DefaultMessageFormat) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	kafkaMsg := message.NewMessage(uuid, msgJSON)
	if err := m.KafkaPublisher.Publish(messageschema.ReplyMessage, kafkaMsg); err != nil {
		return fmt.Errorf("failed on publish message: %v", err)
	}

	return nil
}
