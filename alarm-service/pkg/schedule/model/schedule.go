package model

type WorkSchedule struct {
	RefID   string   `bson:"refId"`
	Owner   string   `bson:"owner"`
	Topic   string   `bson:"topic"`
	Message string   `bson:"message"`
	Time    WorkTime `bson:"time"`
	Type    string   `bson:"type"`
}

type WorkTime struct {
	Date   string `bson:"date"`
	Hour   int    `bson:"hour"`
	Minute int    `bson:"minute"`
	Second int    `bson:"second"`
}
