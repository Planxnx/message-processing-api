package cron

import (
	"context"
	"log"

	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/repository"

	"github.com/robfig/cron/v3"
)

type CronCommand struct {
	CommandFunction func()
	TimeSpec        string
}

type CronUsecase struct {
	ScheduleRepository *repository.ScheduleRepository
	Cron               *cron.Cron
	CronCommands       []CronCommand
}

func NewScheduleUsecase(schRepo *repository.ScheduleRepository) *CronUsecase {
	return &CronUsecase{
		ScheduleRepository: schRepo,
		Cron:               cron.New(),
		CronCommands: []CronCommand{
			{
				CommandFunction: dailyWorkCommand(schRepo),
				TimeSpec:        "@every 1m",
			},
			{
				CommandFunction: hourlyWorkCommand(schRepo),
				TimeSpec:        "@every 1m",
			},
			{
				CommandFunction: weeklyWorkCommand(schRepo),
				TimeSpec:        "@every 1m",
			},
		},
	}
}

func (sch *CronUsecase) StartFetchSchedule() {
	for _, cronCmd := range sch.CronCommands {
		sch.Cron.AddFunc(cronCmd.TimeSpec, cronCmd.CommandFunction)
	}
	sch.Cron.Start()
	log.Println("Start Fetch schedule")
}

func dailyWorkCommand(schR *repository.ScheduleRepository) func() {
	log.Println("Initial DailyWorkCommand schedule")
	return func() {
		ctx := context.Background()

		log.Println("Start daily work schedule!")
		workSch, err := schR.GetDailySchedule(ctx)
		if err != nil {
			log.Println("dailyWorkCommand Error: failed on get daily schedule: " + err.Error())
			return
		}
		if workSch == nil {
			return
		}
		log.Printf("Daily Found! %v", workSch)
		//TODO Publish works to Kafka
	}
}

func hourlyWorkCommand(schR *repository.ScheduleRepository) func() {
	log.Println("Initial HourlyWorkCommand schedule")
	return func() {
		ctx := context.Background()

		log.Println("Start houry work schedule!")
		workSch, err := schR.GetHOURLYSchedule(ctx)
		if err != nil {
			log.Println("hourlyWorkCommand Error: failed on get hour schedule: " + err.Error())
			return
		}
		if workSch == nil {
			return
		}
		log.Printf("Hourly Found! %v", workSch)
		//TODO Publish works to Kafka
	}
}

func weeklyWorkCommand(schR *repository.ScheduleRepository) func() {
	log.Println("Initial WeeklyWorkCommand schedule")
	return func() {
		ctx := context.Background()

		log.Println("Start weekly work schedule!")
		workSch, err := schR.GetWeeklySchedule(ctx)
		if err != nil {
			log.Println("weeklyWorkCommand Error: failed on get weekly schedule: " + err.Error())
			return
		}
		if workSch == nil {
			return
		}
		log.Printf("Weekly Found: %v", workSch)
		//TODO Publish works to Kafka
	}
}
