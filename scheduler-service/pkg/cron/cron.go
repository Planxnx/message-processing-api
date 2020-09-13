package cron

import (
	"context"
	"encoding/json"
	"log"
	"time"

	kafaPkg "github.com/Planxnx/message-processing-api/scheduler-service/pkg/connection/kafka"
	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/repository"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
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
	KafkaPublisher     *kafka.Publisher
}

func NewScheduleUsecase(schRepo *repository.ScheduleRepository, kafkaPubliser *kafka.Publisher) *CronUsecase {
	return &CronUsecase{
		ScheduleRepository: schRepo,
		KafkaPublisher:     kafkaPubliser,
		Cron:               cron.New(),
		CronCommands: []CronCommand{
			{
				CommandFunction: dailyWorkCommand(schRepo, kafkaPubliser),
				TimeSpec:        "@every 1m",
			},
			{
				CommandFunction: hourlyWorkCommand(schRepo, kafkaPubliser),
				TimeSpec:        "@every 1m",
			},
			{
				CommandFunction: weeklyWorkCommand(schRepo, kafkaPubliser),
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

func dailyWorkCommand(schR *repository.ScheduleRepository, kafkaPubliser *kafka.Publisher) func() {
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
		for _, work := range *workSch {
			work.Data["ScheduleType"] = work.Type
			work.Data["ScheduleTime"] = work.Time
			kafkaMessage := &kafaPkg.DefaultMessageFormat{
				Ref1:        work.Ref1,
				Ref2:        work.Ref2,
				Ref3:        work.Ref3,
				Message:     work.Message,
				Owner:       work.Owner,
				PublishedBy: "scheduler-service",
				PublishedAt: time.Now(),
				Features:    work.Features,
				Data:        work.Data,
				Type:        "notification",
			}
			msgJSON, err := json.Marshal(kafkaMessage)
			if err != nil {
				log.Printf("dailyWorkCommand Error: failed on create json message of ref2(%v)and owner(%v): %v ,", work.Ref2, work.Owner, err.Error())
				return
			}
			kafkaMsg := message.NewMessage(watermill.NewUUID(), msgJSON)
			for _, callbackTopic := range work.CallbackTopic {
				if err := kafkaPubliser.Publish(callbackTopic, kafkaMsg); err != nil {
					log.Println("dailyWorkCommand Error: failed on publish message: " + err.Error())
				}
			}
		}
	}
}

func hourlyWorkCommand(schR *repository.ScheduleRepository, kafkaPubliser *kafka.Publisher) func() {
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
		for _, work := range *workSch {
			work.Data["ScheduleType"] = work.Type
			work.Data["ScheduleTime"] = work.Time
			kafkaMessage := &kafaPkg.DefaultMessageFormat{
				Ref1:        work.Ref1,
				Ref2:        work.Ref2,
				Ref3:        work.Ref3,
				Message:     work.Message,
				Owner:       work.Owner,
				PublishedBy: "scheduler-service",
				PublishedAt: time.Now(),
				Features:    work.Features,
				Data:        work.Data,
				Type:        "notification",
			}
			msgJSON, err := json.Marshal(kafkaMessage)
			if err != nil {
				log.Printf("hourlyWorkCommand Error: failed on create json message of ref2(%v)and owner(%v): %v ,", work.Ref2, work.Owner, err.Error())
				return
			}
			kafkaMsg := message.NewMessage(watermill.NewUUID(), msgJSON)
			for _, callbackTopic := range work.CallbackTopic {
				if err := kafkaPubliser.Publish(callbackTopic, kafkaMsg); err != nil {
					log.Println("hourlyWorkCommand Error: failed on publish message: " + err.Error())
				}
			}
		}
	}
}

func weeklyWorkCommand(schR *repository.ScheduleRepository, kafkaPubliser *kafka.Publisher) func() {
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
		for _, work := range *workSch {
			work.Data["ScheduleType"] = work.Type
			work.Data["ScheduleTime"] = work.Time
			kafkaMessage := &kafaPkg.DefaultMessageFormat{
				Ref1:        work.Ref1,
				Ref2:        work.Ref2,
				Ref3:        work.Ref3,
				Message:     work.Message,
				Owner:       work.Owner,
				PublishedBy: "scheduler-service",
				PublishedAt: time.Now(),
				Features:    work.Features,
				Data:        work.Data,
				Type:        "notification",
			}
			msgJSON, err := json.Marshal(kafkaMessage)
			if err != nil {
				log.Printf("weeklyWorkCommand Error: failed on create json message of ref2(%v)and owner(%v): %v ,", work.Ref2, work.Owner, err.Error())
				return
			}
			kafkaMsg := message.NewMessage(watermill.NewUUID(), msgJSON)
			for _, callbackTopic := range work.CallbackTopic {
				if err := kafkaPubliser.Publish(callbackTopic, kafkaMsg); err != nil {
					log.Println("weeklyWorkCommand Error: failed on publish message: " + err.Error())
				}
			}
		}
	}
}
