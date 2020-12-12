package health

import "time"

type HealthData struct {
	Feature       string    `bson:"feature" json:"feature"`
	Description   string    `bson:"description" json:"description"`
	ExecuteMode   []string  `bson:"executeMode" json:"executeMode"`
	ServiceName   string    `bson:"serviceName" json:"serviceName"`
	LastCheckedAt time.Time `bson:"lastCheckedAt" json:"lastCheckedAt"`
}

type HealthLog struct {
	Feature     string    `bson:"feature" json:"feature"`
	Description string    `bson:"description" json:"description"`
	ExecuteMode []string  `bson:"executeMode" json:"executeMode"`
	ServiceName string    `bson:"serviceName" json:"serviceName"`
	CheckedAt   time.Time `bson:"checkedAt" json:"checkedAt"`
}
