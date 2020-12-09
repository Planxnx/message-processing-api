package message

import (
	"encoding/json"
	"log"
	"time"

	botnoiusecase "github.com/Planxnx/message-processing-api/botnoi-service/internal/botnoi"
	messageusecase "github.com/Planxnx/message-processing-api/botnoi-service/internal/message"
	messageSchema "github.com/Planxnx/message-processing-api/message-schema"
	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageHandler struct {
	messageUsecase *messageusecase.MessageUsecase
	botnoiUsecase  *botnoiusecase.BotnoiService
}

func New(messageUsecase *messageusecase.MessageUsecase, botnoiUsecase *botnoiusecase.BotnoiService) *MessageHandler {
	return &MessageHandler{
		messageUsecase: messageUsecase,
		botnoiUsecase:  botnoiUsecase,
	}
}

func (m *MessageHandler) ChitchatHandler(msg *message.Message) error {
	resultMsg := &messageSchema.DefaultMessageFormat{}
	json.Unmarshal(msg.Payload, resultMsg)

	if !resultMsg.Features["chitchat"] {
		return nil
	}

	replyMessage, err := m.botnoiUsecase.ChitChatMessage(resultMsg.Message)
	if err != nil {
		log.Printf("ChitchatHandler Error: failed on chitchat msg: %v", err)
		return err
	}

	if resultMsg.CallbackFlag {
		err = m.messageUsecase.Emit(watermill.NewUUID(), resultMsg.CallbackTopic, &messageschema.DefaultMessageFormat{
			Ref1:        resultMsg.Ref1,
			Ref2:        resultMsg.Ref2,
			Ref3:        resultMsg.Ref3,
			Owner:       "Gateway service",
			PublishedBy: "Botnoi service",
			PublishedAt: time.Now(),
			Data: map[string]interface{}{
				"message": replyMessage,
			},
			Type: "replyMessage",
		})
		if err != nil {
			log.Printf("ChitchatHandler Error: failed on emit message: %v", err)
			return err
		}
		log.Printf("OK(specifi topic)!!!!!!!!")
		return nil
	}

	err = m.messageUsecase.EmitReply(watermill.NewUUID(), &messageschema.DefaultMessageFormat{
		Ref1:        resultMsg.Ref1,
		Ref2:        resultMsg.Ref2,
		Ref3:        resultMsg.Ref3,
		Owner:       "Gateway service",
		PublishedBy: "Botnoi service",
		PublishedAt: time.Now(),
		Data: map[string]interface{}{
			"message": replyMessage,
		},
		Type: "replyMessage",
	})
	if err != nil {
		log.Printf("ChitchatHandler Error: failed on emit message: %v", err)
		return err
	}

	log.Printf("OK!!!!!!!!")
	return nil
}
