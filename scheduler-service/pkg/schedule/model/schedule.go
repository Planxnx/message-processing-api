package model

import "time"

type WorkSchedule struct {
	Ref1          string                 `bson:"ref1" json:"ref1"`   //client reference
	Ref2          string                 `bson:"ref2" json:"ref2"`   //message reference
	Ref3          string                 `bson:"ref3" json:"ref3"`   //end-user reference
	Owner         string                 `bson:"owner" json:"owner"` //service reference
	CallbackTopic []string               `bson:"callbackTopic" json:"callbackTopic"`
	Message       string                 `bson:"message" json:"message"`
	Data          map[string]interface{} `bson:"data" json:"data"` //attachment
	Time          WorkTime               `bson:"time" json:"time"`
	Type          string                 `bson:"type" json:"type"` //Daily, Hourly
	CreateAt      time.Time              `bson:"createdAt" json:"createdAt"`
	DeletedAt     time.Time              `bson:"deletedAt" json:"deletedAt"`
	Features      map[string]bool        `bson:"features" json:"features"` //Feature this message will uses next
}

type WorkTime struct {
	Timestamp time.Time `bson:"timestamp" json:"ref1"` //ISO date
	Day       int       `bson:"day" json:"ref1"`
	WeekDay   string    `bson:"weekDay" json:"ref1"`
	Hour      int       `bson:"hour" json:"ref1"`
	Minute    int       `bson:"minute" json:"ref1"`
	Second    int       `bson:"second" json:"ref1"`
}
