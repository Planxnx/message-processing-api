package message

import (
	lotteryusecase "github.com/Planxnx/message-processing-api/external-caller-service/internal/lottery"

	botnoiusecase "github.com/Planxnx/message-processing-api/external-caller-service/internal/botnoi"
	messageusecase "github.com/Planxnx/message-processing-api/external-caller-service/internal/message"
)

type MessageHandler struct {
	messageUsecase *messageusecase.MessageUsecase
	botnoiUsecase  *botnoiusecase.BotnoiService
	lotteryUsecase *lotteryusecase.LotteryUsecase
}

func New(messageUsecase *messageusecase.MessageUsecase, botnoiUsecase *botnoiusecase.BotnoiService, lotteryUsecase *lotteryusecase.LotteryUsecase) *MessageHandler {
	return &MessageHandler{
		messageUsecase: messageUsecase,
		botnoiUsecase:  botnoiUsecase,
		lotteryUsecase: lotteryUsecase,
	}
}
