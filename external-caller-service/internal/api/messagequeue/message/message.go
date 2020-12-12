package message

import (
	"encoding/json"
	"log"
	"time"

	botnoiusecase "github.com/Planxnx/message-processing-api/external-caller-service/internal/botnoi"
	messageusecase "github.com/Planxnx/message-processing-api/external-caller-service/internal/message"
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
	defer msg.Ack()
	resultMsg := &messageSchema.DefaultMessageFormat{}
	json.Unmarshal(msg.Payload, resultMsg)

	if !resultMsg.Features["chitchat"] {
		return nil
	}

	replyMessage, err := m.botnoiUsecase.ChitChatMessage(resultMsg.Message)
	if err != nil {
		replymessage := &messageschema.DefaultMessageFormat{
			Ref1:        resultMsg.Ref1,
			Ref2:        resultMsg.Ref2,
			Ref3:        resultMsg.Ref3,
			Owner:       resultMsg.Owner,
			PublishedBy: "Botnoi service",
			PublishedAt: time.Now(),
			Type:        "replyMessage",
			Error:       err.Error(),
		}
		if resultMsg.ExcuteMode == messageSchema.SynchronousMode {
			m.messageUsecase.Emit(watermill.NewUUID(), resultMsg.CallbackTopic, replymessage)
		} else {
			m.messageUsecase.EmitReply(watermill.NewUUID(), replymessage)
		}
		log.Printf("ChitchatHandler Error: failed on chitchat msg: %v", err)
		return err
	}
	replymessage := &messageschema.DefaultMessageFormat{
		Ref1:        resultMsg.Ref1,
		Ref2:        resultMsg.Ref2,
		Ref3:        resultMsg.Ref3,
		Owner:       resultMsg.Owner,
		PublishedBy: "Botnoi service",
		PublishedAt: time.Now(),
		Data: map[string]interface{}{
			"message": replyMessage,
		},
		Type: "replyMessage",
	}

	log.Println("Replied !!")
	if resultMsg.ExcuteMode == messageSchema.SynchronousMode {
		err = m.messageUsecase.Emit(watermill.NewUUID(), resultMsg.CallbackTopic, replymessage)
		if err != nil {
			log.Printf("ChitchatHandler Error: failed on emit message: %v", err)
			return err
		}
		return nil
	}
	err = m.messageUsecase.EmitReply(watermill.NewUUID(), replymessage)
	if err != nil {
		log.Printf("ChitchatHandler Error: failed on emit message: %v", err)
		return err
	}

	return nil
}
