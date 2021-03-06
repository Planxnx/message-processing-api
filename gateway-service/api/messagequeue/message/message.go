package message

import (
	"context"
	"encoding/json"
	"log"
	"time"

	callbackusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/callback"
	messageusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/message"
	providerusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/provider"
	"google.golang.org/protobuf/proto"

	messageSchema "github.com/Planxnx/message-processing-api/message-schema"

	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageHandler struct {
	messageUsecase  *messageusecase.MessageUsecase
	providerUsecase *providerusecase.ProviderUsercase
	callbackUsecase *callbackusecase.CallbackUsecase
}

func New(m *messageusecase.MessageUsecase, p *providerusecase.ProviderUsercase, cu *callbackusecase.CallbackUsecase) *MessageHandler {
	return &MessageHandler{
		messageUsecase:  m,
		providerUsecase: p,
		callbackUsecase: cu,
	}
}

func (m *MessageHandler) ReplyMessage(msg *message.Message) error {
	ctx := context.Background()
	resultMsg := &messageSchema.DefaultMessage{}
	proto.Unmarshal(msg.Payload, resultMsg)

	provider, err := m.providerUsecase.GetProviderByID(ctx, resultMsg.Ref1)
	if err != nil {
		log.Printf("ReplyMessage Error: failed on get provider: %v", err)
		return nil
	}

	replyMessage := map[string]interface{}{
		"messageRef": resultMsg.Ref2,
		"ref1":       resultMsg.Ref1,
		"ref2":       resultMsg.Ref2,
		"ref3":       resultMsg.Ref3,
		"type":       resultMsg.Type,
	}

	if resultMsg.ErrorInternal != "" {
		log.Printf("ReplyMessage Error: failed on result: %v", resultMsg.ErrorInternal)
		replyMessage["error"] = "Internal Server Error"
	} else if resultMsg.Error != "" {
		replyMessage["error"] = resultMsg.Error
	} else {
		attachmentData := &map[string]interface{}{}
		json.Unmarshal(resultMsg.Data, attachmentData)
		replyMessage["data"] = attachmentData
	}

	time.Sleep(1 * time.Second) //Add deley
	_, err = m.callbackUsecase.Request(provider.Webhook, replyMessage)
	if err != nil {
		log.Printf("ReplyMessage Error: failed on send callback to webhook: %v", err)
		return nil
	}
	log.Println("==================")
	return nil
}
