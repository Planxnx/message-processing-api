package message

import (
	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/google/uuid"
)

type MessageUsecase struct{}

func New() *MessageUsecase {
	return &MessageUsecase{}
}

func (MessageUsecase) Emit(msg *messageschema.DefaultMessageFormat) (string, error) {
	return uuid.New().String(), nil
}
