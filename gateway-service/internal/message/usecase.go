package message

import (
	"errors"
	"fmt"
	"strings"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type Config struct {
	FeatureName      string
	ServiceName      string
	AsynchronousMode bool
	SynchronousMode  bool
}

func (m *MessageUsecase) CreateMessageHandler(config *Config, handler func(response, request *messageschema.DefaultMessage) error) func(msg *message.Message) error {

	return func(msg *message.Message) error {
		defer msg.Ack()

		requestmessage := &messageschema.DefaultMessage{}
		proto.Unmarshal(msg.Payload, requestmessage)

		//Validate feture name
		if strings.ToLower(requestmessage.Feature) != strings.ToLower(config.FeatureName) {
			return nil
		}

		//Define default reply message
		replymessage := &messageschema.DefaultMessage{
			Ref1:          requestmessage.Ref1,
			Ref2:          requestmessage.Ref2,
			Ref3:          requestmessage.Ref3,
			Owner:         requestmessage.Owner,
			Feature:       requestmessage.Feature,
			ExcuteMode:    requestmessage.ExcuteMode,
			CallbackTopic: requestmessage.CallbackTopic,
			PublishedBy:   config.ServiceName,
			Type:          "replyMessage",
		}

		//Reply Handler
		defer func() {
			//Reply for Synchronous
			if requestmessage.ExcuteMode == messageschema.ExecuteMode_Synchronous {
				m.Emit(watermill.NewUUID(), requestmessage.CallbackTopic, replymessage)
				return
			}

			//Reply for Asynchronous
			m.EmitReply(watermill.NewUUID(), replymessage)
			return
		}()

		//validate support mode
		if requestmessage.ExcuteMode != messageschema.ExecuteMode_Synchronous && requestmessage.ExcuteMode != messageschema.ExecuteMode_Asynchronous {
			return nil
		}

		if (!config.AsynchronousMode && requestmessage.ExcuteMode == messageschema.ExecuteMode_Asynchronous) || (!config.SynchronousMode && requestmessage.ExcuteMode == messageschema.ExecuteMode_Synchronous) {
			replymessage.Error = "wrong exec mode"
			replymessage.PublishedAt = timestamppb.Now()
			return errors.New(replymessage.Error)
		}

		return handler(replymessage, requestmessage)
	}

}
