package message

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	kafapkg "github.com/Planxnx/message-processing-api/gateway-service/pkg/kafka"
	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/gofiber/fiber/v2"

	healthusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/health"
	messageusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/message"

	"github.com/Planxnx/message-processing-api/gateway-service/model"
	messageschema "github.com/Planxnx/message-processing-api/message-schema"
)

type MessageHandler struct {
	MessageUsecase  *messageusecase.MessageUsecase
	KafkaSubscriber *kafka.Subscriber
	healthUsercase  *healthusecase.HealthUsercase
	SaramaBroker    *sarama.Broker
}

func New(m *messageusecase.MessageUsecase, sub *kafka.Subscriber, hu *healthusecase.HealthUsercase, sb *sarama.Broker) *MessageHandler {
	return &MessageHandler{
		MessageUsecase:  m,
		KafkaSubscriber: sub,
		healthUsercase:  hu,
		SaramaBroker:    sb,
	}
}

//MainEndpoint Asyncmode
func (m *MessageHandler) MainEndpoint(c *fiber.Ctx) error {
	providerID := c.Get("Provider-ID")
	reqBody := &model.MessageRequest{}
	c.BodyParser(reqBody)

	featureHealth, _ := m.healthUsercase.GetHealth(c.Context(), reqBody.Feature)
	if featureHealth == nil {
		return c.Status(fiber.StatusBadRequest).JSON(&model.Response{
			Message: "feature is unavailable",
		})
	}
	index := sort.SearchStrings(featureHealth.ExecuteMode, messageschema.ExecuteMode_Asynchronous.String())
	if featureHealth.ExecuteMode[index] != messageschema.ExecuteMode_Asynchronous.String() {
		return c.Status(fiber.StatusBadRequest).JSON(&model.Response{
			Message: "asynchronous mode not support",
		})
	}

	messageRef := watermill.NewUUID()

	dataByte, err := json.Marshal(reqBody.Data)
	if err != nil {
		log.Printf("MainEndpoint Error: failed on marshal req data: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error or request.data is invalid format")
	}

	err = m.MessageUsecase.EmitCommon(messageRef, &messageschema.DefaultMessage{
		Message:     reqBody.Message,
		Ref1:        providerID,
		Ref2:        messageRef,
		Ref3:        reqBody.UserRef,
		Owner:       "Gateway service",
		PublishedBy: "Gateway service",
		PublishedAt: timestamppb.Now(),
		Feature:     reqBody.Feature,
		Data:        dataByte,
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
		log.Fatalf("SynchronousEndpoint: failed on create kafka subscriber: %v", err)
	}
	defer kafkaSubscriber.Close()

	providerID := c.Get("Provider-ID")
	reqBody := &model.MessageRequest{}
	c.BodyParser(reqBody)

	featureHealth, _ := m.healthUsercase.GetHealth(c.Context(), reqBody.Feature)
	if featureHealth == nil {
		return c.Status(fiber.StatusBadRequest).JSON(&model.Response{
			Message: "feature is unavailable",
		})
	}
	index := sort.SearchStrings(featureHealth.ExecuteMode, messageschema.ExecuteMode_Synchronous.String())
	if featureHealth.ExecuteMode[index] != messageschema.ExecuteMode_Synchronous.String() {
		return c.Status(fiber.StatusBadRequest).JSON(&model.Response{
			Message: "synchronous mode not support",
		})
	}

	messageRef := watermill.NewShortUUID()
	callbackTopic := fmt.Sprintf("response-%v", messageRef)

	dataByte, _ := json.Marshal(reqBody.Data)

	err = m.MessageUsecase.EmitCommon(messageRef, &messageschema.DefaultMessage{
		Message:       reqBody.Message,
		Ref1:          providerID,
		Ref2:          messageRef,
		Ref3:          reqBody.UserRef,
		Owner:         "Gateway service",
		PublishedBy:   "Gateway service",
		PublishedAt:   timestamppb.Now(),
		Feature:       reqBody.Feature,
		Data:          dataByte,
		Type:          "newMessage",
		ExcuteMode:    messageschema.ExecuteMode_Synchronous,
		CallbackTopic: callbackTopic,
	})
	if err != nil {
		log.Printf("SynchronousEndpoint Error: failed on emit message: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	if err := kafapkg.CreateTopic(callbackTopic, m.SaramaBroker); err != nil {
		log.Printf("SynchronousEndpoint Error: failed on create kafka topic: %v\n", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}
	defer kafapkg.DeleteTopic([]string{callbackTopic}, m.SaramaBroker)

	submessage, err := kafkaSubscriber.Subscribe(ctx, callbackTopic)
	if err != nil {
		log.Printf("SynchronousEndpoint Error: failed on subscribe message: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	respmessage := <-submessage
	respmessage.Ack()

	resultMsg := &messageschema.DefaultMessage{}
	proto.Unmarshal(respmessage.Payload, resultMsg)

	if resultMsg.Error != "" {
		return fiber.NewError(fiber.StatusBadRequest, resultMsg.Error)
	}

	if resultMsg.ErrorInternal != "" {
		log.Printf("SynchronousEndpoint Error: failed on result: %v", resultMsg.ErrorInternal)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	attachmentData := &map[string]interface{}{}
	json.Unmarshal(resultMsg.Data, attachmentData)

	return c.Status(fiber.StatusOK).JSON(&model.Response{
		Message: "success",
		Data:    attachmentData,
	})
}
