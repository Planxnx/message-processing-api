package message

import (
	"encoding/json"
	"log"
	"strings"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
		replymessage.ErrorInternal = err.Error()
		replymessage.PublishedAt = timestamppb.Now()
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
		replymessage.ErrorInternal = err.Error()
		replymessage.PublishedAt = timestamppb.Now()
		if resultMsg.ExcuteMode == messageschema.ExecuteMode_Synchronous {
			m.messageUsecase.Emit(watermill.NewUUID(), resultMsg.CallbackTopic, replymessage)
		} else {
			m.messageUsecase.EmitReply(watermill.NewUUID(), replymessage)
		}
		log.Printf("ChitchatHandler Error: failed on marshal attchment: %v", err)
		return err
	}

	replymessage.Data = attachmentData
	replymessage.PublishedAt = timestamppb.Now()
	replymessage.PublishedAt = timestamppb.Now()

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
