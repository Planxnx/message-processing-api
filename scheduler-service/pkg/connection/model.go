package connection

import "time"

type WorkSchedule struct {
	Owner   string
	Message string
	Time    time.Time
}
