package messageapi

import (
	"context"
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
	Handler MessageAPIHandler
	Config  *HandlerConfig
}

type Handler func(msg *message.Message) error

type MessageAPIHandler func(response, request *messageschema.DefaultMessage) error

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

//AddHandler add
func (m *MessageClient) AddHandler(topic string, config *HandlerConfig, fn func(response, request *messageschema.DefaultMessage) error) {
	m.handlers = append(m.handlers, &handler{
		Topic:   topic,
		Handler: fn,
		Config:  config,
	})
}

//CreateMessageHandler create message api handler
func (m *MessageClient) CreateMessageAPIHandler(config *HandlerConfig, handler MessageAPIHandler) Handler {

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
				m.CreateMessageAPIHandler(h.Config, h.Handler)(msg)
			}
		}(v)

		go m.healthCheck(v)
	}

	select {
	case <-ctx.Done():
		return errors.Errorf("ctx cancelled")
	}
}
