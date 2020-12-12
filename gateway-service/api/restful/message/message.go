package message

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	kafapkg "github.com/Planxnx/message-processing-api/gateway-service/pkg/kafka"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/gofiber/fiber/v2"

	messageusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/message"
	"github.com/Planxnx/message-processing-api/gateway-service/model"
	messageschema "github.com/Planxnx/message-processing-api/message-schema"
)

type MessageHandler struct {
	MessageUsecase  *messageusecase.MessageUsecase
	KafkaSubscriber *kafka.Subscriber
}

func New(m *messageusecase.MessageUsecase, sub *kafka.Subscriber) *MessageHandler {
	return &MessageHandler{
		MessageUsecase:  m,
		KafkaSubscriber: sub,
	}
}

func (m *MessageHandler) MainEndpoint(c *fiber.Ctx) error {
	providerID := c.Get("Provider-ID")
	reqBody := &model.MessageRequest{}
	c.BodyParser(reqBody)
	messageRef := watermill.NewUUID()
	err := m.MessageUsecase.EmitCommon(messageRef, &messageschema.DefaultMessageFormat{
		Message:     reqBody.Message,
		Ref1:        providerID,
		Ref2:        messageRef,
		Ref3:        reqBody.UserRef,
		Owner:       "Gateway service",
		PublishedBy: "Gateway service",
		PublishedAt: time.Now(),
		Features:    reqBody.Features,
		Data:        reqBody.Data,
		Type:        "newMessage",
	})
	if err != nil {
		log.Printf("MainEndpoint Error: failed on emit message: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}
	return c.Status(fiber.StatusOK).JSON(&model.Response{
		Message: "Success",
		Data: model.MessageResponseData{
			MessageRef: messageRef,
		},
	})
}

func (m *MessageHandler) SynchronousEndpoint(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	kafkaSubscriber, err := kafapkg.NewSubscriber()
	if err != nil {
		log.Fatalf("main Error: failed on create kafka subscriber: %v", err)
	}
	defer kafkaSubscriber.Close()

	providerID := c.Get("Provider-ID")
	reqBody := &model.MessageRequest{}
	c.BodyParser(reqBody)
	messageRef := watermill.NewUUID()
	callbackTopic := fmt.Sprintf("response-%v", messageRef)

	submessage, err := kafkaSubscriber.Subscribe(ctx, callbackTopic)
	if err != nil {
		log.Printf("MainEndpoint Error: failed on subscribe message: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	err = m.MessageUsecase.EmitCommon(messageRef, &messageschema.DefaultMessageFormat{
		Message:       reqBody.Message,
		Ref1:          providerID,
		Ref2:          messageRef,
		Ref3:          reqBody.UserRef,
		Owner:         "Gateway service",
		PublishedBy:   "Gateway service",
		PublishedAt:   time.Now(),
		Features:      reqBody.Features,
		Data:          reqBody.Data,
		Type:          "newMessage",
		ExcuteMode:    messageschema.SynchronousMode,
		CallbackTopic: callbackTopic,
	})
	if err != nil {
		log.Printf("MainEndpoint Error: failed on emit message: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	respmessage := <-submessage
	respmessage.Ack()

	resultMsg := &messageschema.DefaultMessageFormat{}
	json.Unmarshal(respmessage.Payload, resultMsg)

	if resultMsg.Error != "" {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	return c.Status(fiber.StatusOK).JSON(&model.Response{
		Message: "success",
		Data:    resultMsg.Data,
	})
}
