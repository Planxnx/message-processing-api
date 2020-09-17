package message

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/Planxnx/message-processing-api/gateway-service/model"

	messageusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/message"
	messageschema "github.com/Planxnx/message-processing-api/message-schema"
)

type MessageHandler struct {
	MessageUsecase *messageusecase.MessageUsecase
}

func New(m *messageusecase.MessageUsecase) *MessageHandler {
	return &MessageHandler{
		MessageUsecase: m,
	}
}

func (m *MessageHandler) MainEndpoint(c *fiber.Ctx) error {
	reqBody := &model.MessageRequest{}
	c.BodyParser(reqBody)
	messageRef := uuid.New().String()
	mockID, err := m.MessageUsecase.Emit(&messageschema.DefaultMessageFormat{
		Message:     reqBody.Message,
		Ref1:        c.IP(),
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
	return c.Status(fiber.StatusOK).JSON(&model.MessageResponse{
		Message: "Success",
		Data: model.MessageResponseData{
			MessageRef: mockID,
		},
	})
}