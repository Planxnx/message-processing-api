package message

import (
	"encoding/json"
	"log"

	messageusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/message"
	messageSchema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageHandler struct {
	MessageUsecase *messageusecase.MessageUsecase
}

func New(m *messageusecase.MessageUsecase) *MessageHandler {
	return &MessageHandler{
		MessageUsecase: m,
	}
}

func (m *MessageHandler) ReplyMessage(msg *message.Message) error {
	resultMsg := &messageSchema.DefaultMessageFormat{}
	json.Unmarshal(msg.Payload, resultMsg)

	//TODO: send to webhook
	log.Printf("Received Notification Event:\n  ---Unmarshal Message: %v \n  ---Raw Message: %v \n", resultMsg, string(msg.Payload))

	return nil
}
