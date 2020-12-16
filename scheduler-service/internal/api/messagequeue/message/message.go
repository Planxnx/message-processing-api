package message

import (
	"github.com/Planxnx/message-processing-api/scheduler-service/internal/alarm"
	messageusecase "github.com/Planxnx/message-processing-api/scheduler-service/internal/message"
)

type MessageHandler struct {
	messageUsecase *messageusecase.MessageUsecase
	alarmService   *alarm.Service
}

func New(messageUsecase *messageusecase.MessageUsecase, alarmService *alarm.Service) *MessageHandler {
	return &MessageHandler{
		messageUsecase: messageUsecase,
		alarmService:   alarmService,
	}
}
