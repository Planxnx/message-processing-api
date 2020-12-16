package model

import (
	"time"

	scheduleconstant "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/constant"
)

type WorkSchedule struct {
	Ref1          string                        `bson:"ref1" json:"ref1"`   //client reference
	Ref2          string                        `bson:"ref2" json:"ref2"`   //message reference
	Ref3          string                        `bson:"ref3" json:"ref3"`   //end-user reference
	Owner         string                        `bson:"owner" json:"owner"` //service reference
	CallbackTopic string                        `bson:"callbackTopic" json:"callbackTopic"`
	Message       string                        `bson:"message" json:"message"`
	Data          map[string]interface{}        `bson:"data" json:"data"` //attachment
	Time          WorkTime                      `bson:"time" json:"time"`
	Type          scheduleconstant.ScheduleType `bson:"type" json:"type"` //Daily, Hourly
	CreateAt      time.Time                     `bson:"createdAt" json:"createdAt"`
	DeletedAt     time.Time                     `bson:"deletedAt" json:"deletedAt"`
	Feature       string                        `bson:"feature" json:"feature"` //Feature this message will uses next
	IsOnce        bool                          `bson:"isOnce" json:"isOnce"`   //is one time usage ?
}

type WorkTime struct {
	Timestamp time.Time                `bson:"timestamp" json:"timestamp"` //date
	Day       int                      `bson:"day" json:"day"`
	WeekDay   scheduleconstant.WeekDay `bson:"weekDay" json:"weekDay"` //required for weekly
	Hour      int                      `bson:"hour" json:"hour"`       //required for weekly,daily
	Minute    int                      `bson:"minute" json:"minute"`   //require for weekly,daily,hourly
	Second    int                      `bson:"second" json:"second"`
}
