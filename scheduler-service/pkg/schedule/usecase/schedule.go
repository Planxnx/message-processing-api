package schedule

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

type ScheduleUsecase struct {
	ScheduleRepository *repository.ScheduleRepository
	Cron               *cron.Cron
	CronCommands       []CronCommand
}

func NewScheduleUsecase(schRepo *repository.ScheduleRepository) *ScheduleUsecase {
	return &ScheduleUsecase{
		ScheduleRepository: schRepo,
		Cron:               cron.New(),
		CronCommands: []CronCommand{
			{
				CommandFunction: dailyWorkCommand(schRepo),
				TimeSpec:        "@every 1m",
			},
			{
				CommandFunction: everyHourWorkCommand(schRepo),
				TimeSpec:        "@every 1m",
			},
		},
	}
}

func (sch *ScheduleUsecase) StartFetchSchedule() {
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
		log.Printf("Found! %v", workSch)
		//TODO Publish works to Kafka
	}
}

func everyHourWorkCommand(schR *repository.ScheduleRepository) func() {
	log.Println("Initial EveryHourWorkCommand schedule")
	return func() {
		ctx := context.Background()

		log.Println("Start everyhour work schedule!")
		workSch, err := schR.GetEveryHourSchedule(ctx)
		if err != nil {
			log.Println("everyHourWorkCommand Error: failed on get hour schedule: " + err.Error())
			return
		}
		if workSch == nil {
			return
		}

		//TODO Publish works to Kafka
	}
}
