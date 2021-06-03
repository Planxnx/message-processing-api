package messageapi

import (
	"fmt"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/protobuf/proto"
)

func (m *MessageClient) EmitCommon(uuid string, msg *messageschema.DefaultMessage) error {
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

func (m *MessageClient) EmitReply(uuid string, msg *messageschema.DefaultMessage) error {
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

func (m *MessageClient) Emit(uuid string, topic string, msg *messageschema.DefaultMessage) error {
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
