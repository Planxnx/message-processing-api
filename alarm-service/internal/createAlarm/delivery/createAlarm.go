package delivery

import (
	"encoding/json"
	"log"

	messageSchema "github.com/Planxnx/message-processing-api/message-schema"

	createAlarmUsecase "github.com/Planxnx/message-processing-api/alarm-service/internal/createAlarm/usecase"
	"github.com/ThreeDotsLabs/watermill/message"
)

type CreateAlarmDelivery struct {
	CreateAlarmUsecase *createAlarmUsecase.CreateAlarmUsecase
}

func NewCreateAlarmDelivery(caUsecase *createAlarmUsecase.CreateAlarmUsecase) *CreateAlarmDelivery {
	return &CreateAlarmDelivery{
		CreateAlarmUsecase: caUsecase,
	}
}

func (ca *CreateAlarmDelivery) CommonMessageHandler(msg *message.Message) error {
	//TODO: implement commond msg
	ctx := msg.Context()
	resultMsg := &messageSchema.DefaultMessageFormat{}
	json.Unmarshal(msg.Payload, resultMsg)

	if resultMsg.Features["dailyAlarm"] {
		_ = ca.CreateAlarmUsecase.CreateDailyAlarm(ctx)
	}
	log.Printf("Received Common Event:\n  ---Unmarshal Message: %v \n  ---Raw Message: %v \n", resultMsg, string(msg.Payload))
	return nil
}
