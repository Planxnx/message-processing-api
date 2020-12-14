package message

import (
	"encoding/json"
	"log"
	"strings"

	botnoiusecase "github.com/Planxnx/message-processing-api/external-caller-service/internal/botnoi"
	messageusecase "github.com/Planxnx/message-processing-api/external-caller-service/internal/message"
	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
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
	resultMsg := &messageschema.DefaultMessage{}
	proto.Unmarshal(msg.Payload, resultMsg)

	if strings.ToLower(resultMsg.Feature) != "chitchat" {
		return nil
	}

	//validate support mode
	if resultMsg.ExcuteMode != messageschema.ExecuteMode_Synchronous && resultMsg.ExcuteMode != messageschema.ExecuteMode_Asynchronous {
		log.Println("Wrong ExecMode")
		return nil
	}

	replymessage := &messageschema.DefaultMessage{
		Ref1:        resultMsg.Ref1,
		Ref2:        resultMsg.Ref2,
		Ref3:        resultMsg.Ref3,
		Owner:       resultMsg.Owner,
		PublishedBy: "External Caller service",
		Type:        "replyMessage",
	}

	replyMessage, err := m.botnoiUsecase.ChitChatMessage(resultMsg.Message)
	if err != nil {
		replymessage.Error = err.Error()
		replymessage.PublishedAt = ptypes.TimestampNow()
		if resultMsg.ExcuteMode == messageschema.ExecuteMode_Synchronous {
			m.messageUsecase.Emit(watermill.NewUUID(), resultMsg.CallbackTopic, replymessage)
		} else {
			m.messageUsecase.EmitReply(watermill.NewUUID(), replymessage)
		}
		log.Printf("ChitchatHandler Error: failed on chitchat msg: %v", err)
		return err
	}

	attachmentData, err := json.Marshal(map[string]interface{}{
		"message": replyMessage,
	})
	if err != nil {
		replymessage.Error = err.Error()
		replymessage.PublishedAt = ptypes.TimestampNow()
		if resultMsg.ExcuteMode == messageschema.ExecuteMode_Synchronous {
			m.messageUsecase.Emit(watermill.NewUUID(), resultMsg.CallbackTopic, replymessage)
		} else {
			m.messageUsecase.EmitReply(watermill.NewUUID(), replymessage)
		}
		log.Printf("ChitchatHandler Error: failed on marshal attchment: %v", err)
		return err
	}

	replymessage.Data = attachmentData
	replymessage.PublishedAt = ptypes.TimestampNow()
	replymessage.PublishedAt = ptypes.TimestampNow()

	log.Println("Replied !!")
	if resultMsg.ExcuteMode == messageschema.ExecuteMode_Synchronous {
		err = m.messageUsecase.Emit(watermill.NewShortUUID(), resultMsg.CallbackTopic, replymessage)
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
