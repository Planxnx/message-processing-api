package message

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/Planxnx/message-processing-api/scheduler-service/internal/alarm"
	scheduleconstant "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/constant"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const featureName string = "daily-notification"

type DailyNotificationRequestData struct {
	Time       time.Time              `json:"time,omitempty"`
	Attachment map[string]interface{} `json:"attachment,omitempty"`
}

func (m *MessageHandler) RegisterDailyNotificationHandler(msg *message.Message) error {
	ctx := msg.Context()
	defer msg.Ack()
	resultMsg := &messageschema.DefaultMessage{}
	proto.Unmarshal(msg.Payload, resultMsg)

	if strings.ToLower(resultMsg.Feature) != featureName {
		return nil
	}

	//validate support mode
	if resultMsg.ExcuteMode != messageschema.ExecuteMode_Asynchronous {
		log.Println("Wrong ExecMode")
		return nil
	}

	replymessage := &messageschema.DefaultMessage{
		Ref1:        resultMsg.Ref1,
		Ref2:        resultMsg.Ref2,
		Ref3:        resultMsg.Ref3,
		Owner:       resultMsg.Owner,
		PublishedBy: "Scheduler service",
		Type:        "replyMessage",
	}

	requestData := &DailyNotificationRequestData{}
	err := json.Unmarshal(resultMsg.Data, requestData)
	if err != nil {
		replymessage.ErrorInternal = err.Error()
		replymessage.PublishedAt = timestamppb.Now()
		m.messageUsecase.Emit(watermill.NewUUID(), resultMsg.CallbackTopic, replymessage)
		log.Printf("RegisterDailyNotificationHandler Error: failed on unmarshal data: %v", err)
		return err
	}

	requestData.Time = requestData.Time.Local()
	hh, mm, ss := requestData.Time.Clock()

	err = m.alarmService.CreateDailyAlarm(ctx, &alarm.AlarmData{
		Ref1:    resultMsg.Ref1,
		Ref2:    resultMsg.Ref2,
		Ref3:    resultMsg.Ref3,
		Owner:   resultMsg.Owner,
		Message: resultMsg.Message,
		Data:    requestData.Attachment,
		Time: alarm.AlarmDataWorkTime{
			Timestamp: requestData.Time,
			Day:       requestData.Time.Day(),
			WeekDay:   scheduleconstant.WeekDay(requestData.Time.Weekday().String()),
			Hour:      hh,
			Minute:    mm,
			Second:    ss,
		},
		Type:          resultMsg.Type,
		Feature:       "replyMessage",
		CallbackTopic: "replyMessage",
		IsOnce:        true,
	})
	if err != nil {
		replymessage.ErrorInternal = err.Error()
		replymessage.PublishedAt = timestamppb.Now()
		m.messageUsecase.EmitReply(watermill.NewUUID(), replymessage)
		log.Printf("RegisterNotiLotteryHandler Error: failed on create daily alarm: %v", err)
		return err
	}

	return nil

}
