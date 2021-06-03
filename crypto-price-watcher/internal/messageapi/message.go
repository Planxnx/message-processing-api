package messageapi

import (
	"context"
	"fmt"
	"strings"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type handler struct {
	Topic   string
	Handler func(response, request *messageschema.DefaultMessage) error
	Config  *HandlerConfig
}

type HandlerConfig struct {
	FeatureName      string
	ServiceName      string
	Description      string
	AsynchronousMode bool
	SynchronousMode  bool
}

type MessageClient struct {
	KafkaPublisher  *kafka.Publisher
	KafkaSubscriber *kafka.Subscriber
	handlers        []*handler
}

func New() (*MessageClient, error) {

	kafkaSubscriber, err := NewSubscriber()
	if err != nil {
		return nil, errors.Errorf("main Error: failed on create kafka subscriber: %v", err)
	}
	kafkaNewPublisher, err := NewPubliser()
	if err != nil {
		return nil, errors.Errorf("main Error: failed on create kafka publisher: %v", err)
	}

	return &MessageClient{
		KafkaPublisher:  kafkaNewPublisher,
		KafkaSubscriber: kafkaSubscriber,
	}, nil
}

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

func (m *MessageClient) CreateMessageHandler(config *HandlerConfig, handler func(response, request *messageschema.DefaultMessage) error) func(msg *message.Message) error {

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

func (m *MessageClient) AddHandler(topic string, config *HandlerConfig, fn func(response, request *messageschema.DefaultMessage) error) {
	m.handlers = append(m.handlers, &handler{
		Topic:   topic,
		Handler: fn,
		Config:  config,
	})
}

//Run start subscriber
func (m *MessageClient) Run(ctx context.Context) error {
	for _, v := range m.handlers {
		messages, err := m.KafkaSubscriber.Subscribe(ctx, v.Topic)
		if err != nil {
			return errors.Errorf("Error: failed on subscribe topic: %v", err)
		}

		go func(h *handler) {
			for msg := range messages {
				//TODO: Handle Errors
				m.CreateMessageHandler(h.Config, h.Handler)(msg)
			}
		}(v)

		go m.healthCheck(v)
	}

	select {
	case <-ctx.Done():
		return errors.Errorf("ctx cancelled")
	}
}
