package schedule

import (
	"log"

	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/connection"
	"github.com/robfig/cron/v3"
)

type Schedule struct {
	Connection *connection.Connections
	Command    func()
	Cron       *cron.Cron
	CronSpec   string
}

func InitializeSchedule(con *connection.Connections) *Schedule {
	log.Println("Initialize Schedule")

	return &Schedule{
		Connection: con,
		Cron:       cron.New(),
		CronSpec:   "@every 1m",
	}
}

func (sch *Schedule) StartSchedule() {
	sch.Cron.AddFunc(sch.CronSpec, func() {
		log.Println("Fetching Schedule!")
		workSch, _ := sch.Connection.MockCon.GetNowSchedule()
		log.Println(workSch)
	})
	sch.Cron.Start()
	log.Println("Start Schedule")
}
