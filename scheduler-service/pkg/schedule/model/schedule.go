package model

import "time"

type WorkSchedule struct {
	Ref1          string    `bson:"ref1"`  //client reference
	Ref2          string    `bson:"ref2"`  //message reference
	Ref3          string    `bson:"ref3"`  //end-user reference
	Owner         string    `bson:"owner"` //service reference
	CallbackTopic string    `bson:"callbackTopic"`
	Message       string    `bson:"message"`
	Time          WorkTime  `bson:"time"`
	Type          string    `bson:"type"` //Daily, Hourly
	CreateAt      time.Time `bson:"createdAt"`
	DeletedAt     time.Time `bson:"deleted"`
}

type WorkTime struct {
	Timestamp time.Time `bson:"timestamp"` //ISO date
	Day       int       `bson:"day"`
	WeekDay   string    `bson:"weekDay"`
	Hour      int       `bson:"hour"`
	Minute    int       `bson:"minute"`
	Second    int       `bson:"second"`
}
