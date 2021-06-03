package alarm

import (
	"context"
	"log"
	"time"

	scheduleconstant "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/constant"

	schedulemodel "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/model"
)

type AlarmData struct {
	Ref1          string                 `bson:"ref1" json:"ref1"`   //client reference
	Ref2          string                 `bson:"ref2" json:"ref2"`   //message reference
	Ref3          string                 `bson:"ref3" json:"ref3"`   //end-user reference
	Owner         string                 `bson:"owner" json:"owner"` //service reference
	Message       string                 `bson:"message" json:"message"`
	Data          map[string]interface{} `bson:"data" json:"data"` //attachment
	Time          AlarmDataWorkTime      `bson:"time" json:"time"`
	Type          string                 `bson:"type" json:"type"`       //Daily, Hourly
	Feature       string                 `bson:"feature" json:"feature"` //Feature this message will uses next
	CallbackTopic string                 `bson:"callbackTopic" json:"callbackTopic"`
	IsOnce        bool                   `bson:"isOnce" json:"isOnce"` //is one time usage ?
	CreateAt      time.Time              `bson:"createdAt" json:"createdAt"`
	DeletedAt     time.Time              `bson:"deletedAt" json:"deletedAt"`
}

type AlarmDataWorkTime struct {
	Timestamp time.Time                `bson:"timestamp" json:"timestamp"` //date time type
	Day       int                      `bson:"day" json:"day"`
	WeekDay   scheduleconstant.WeekDay `bson:"weekDay" json:"weekDay"` //required for weekly
	Hour      int                      `bson:"hour" json:"hour"`       //required for weekly,daily
	Minute    int                      `bson:"minute" json:"minute"`   //require for weekly,daily,hourly
	Second    int                      `bson:"second" json:"second"`
}

func (s *Service) CreateDailyAlarm(ctx context.Context, alarmData *AlarmData) error {
	result, err := s.scheduleRepository.InsertWeeklySchedule(ctx, schedulemodel.WorkSchedule{
		Ref1:          alarmData.Ref1,
		Ref2:          alarmData.Ref2,
		Ref3:          alarmData.Ref3,
		Owner:         alarmData.Owner,
		Message:       alarmData.Message,
		CallbackTopic: alarmData.CallbackTopic,
		Time: schedulemodel.WorkTime{
			Timestamp: alarmData.Time.Timestamp,
			Day:       alarmData.Time.Day,
			WeekDay:   alarmData.Time.WeekDay,
			Hour:      alarmData.Time.Hour,
			Minute:    alarmData.Time.Minute,
			Second:    alarmData.Time.Second,
		},
		Data:    alarmData.Data,
		Feature: alarmData.Feature,
		IsOnce:  alarmData.IsOnce,
	})
	if err != nil {
		log.Printf("Error: failed on insert: %s", err.Error())
		return err
	}
	log.Printf("Insert successful: %v", result.InsertedID)
	return nil
}
