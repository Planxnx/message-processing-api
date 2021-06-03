package cron

import (
	"context"
	"encoding/json"
	"log"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/model"
	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/repository"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/robfig/cron/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (sch *CronUsecase) StartFetchSchedule() *cron.Cron {
	for _, cronCmd := range sch.CronCommands {
		sch.Cron.AddFunc(cronCmd.TimeSpec, cronCmd.CommandFunction)
	}
	sch.Cron.Start()
	log.Println("Start Fetch schedule")

	return sch.Cron
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
			go func(work model.WorkSchedule) {
				work.Data["scheduleType"] = work.Type
				work.Data["scheduleTime"] = work.Time
				workData, err := json.Marshal(work.Data)
				if err != nil {
					log.Printf("dailyWorkCommand Error: failed on json work data marshal of ref2(%v)and owner(%v): %v ,", work.Ref2, work.Owner, err.Error())
					return
				}
				kafkaMessage := &messageschema.DefaultMessage{
					Ref1:        work.Ref1,
					Ref2:        work.Ref2,
					Ref3:        work.Ref3,
					Message:     work.Message,
					Owner:       work.Owner,
					PublishedBy: "Scheduler-service",
					PublishedAt: timestamppb.Now(),
					Feature:     work.Feature,
					Data:        workData,
					Type:        "notification",
				}
				msgByte, err := proto.Marshal(kafkaMessage)
				if err != nil {
					log.Printf("dailyWorkCommand Error: failed on create proto message of ref2(%v)and owner(%v): %v ,", work.Ref2, work.Owner, err.Error())
					return
				}
				kafkaMsg := message.NewMessage(watermill.NewUUID(), msgByte)

				if err := kafkaPubliser.Publish(work.CallbackTopic, kafkaMsg); err != nil {
					log.Println("dailyWorkCommand Error: failed on publish message: " + err.Error())
				}
			}(work)
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
			work.Data["scheduleType"] = work.Type
			work.Data["scheduleTime"] = work.Time
			workData, err := json.Marshal(work.Data)
			if err != nil {
				log.Printf("hourlyWorkCommand Error: failed on json work data marshal of ref2(%v)and owner(%v): %v ,", work.Ref2, work.Owner, err.Error())
				return
			}
			kafkaMessage := &messageschema.DefaultMessage{
				Ref1:        work.Ref1,
				Ref2:        work.Ref2,
				Ref3:        work.Ref3,
				Message:     work.Message,
				Owner:       work.Owner,
				PublishedBy: "Scheduler-service",
				PublishedAt: timestamppb.Now(),
				Feature:     work.Feature,
				Data:        workData,
				Type:        "notification",
			}
			msgByte, err := proto.Marshal(kafkaMessage)
			if err != nil {
				log.Printf("hourlyWorkCommand Error: failed on create json message of ref2(%v)and owner(%v): %v ,", work.Ref2, work.Owner, err.Error())
				return
			}
			kafkaMsg := message.NewMessage(watermill.NewUUID(), msgByte)
			if err := kafkaPubliser.Publish(work.CallbackTopic, kafkaMsg); err != nil {
				log.Println("hourlyWorkCommand Error: failed on publish message: " + err.Error())
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
			work.Data["scheduleType"] = work.Type
			work.Data["scheduleTime"] = work.Time
			workData, err := json.Marshal(work.Data)
			if err != nil {
				log.Printf("weeklyWorkCommand Error: failed on json work data marshal of ref2(%v)and owner(%v): %v ,", work.Ref2, work.Owner, err.Error())
				return
			}
			kafkaMessage := &messageschema.DefaultMessage{
				Ref1:        work.Ref1,
				Ref2:        work.Ref2,
				Ref3:        work.Ref3,
				Message:     work.Message,
				Owner:       work.Owner,
				PublishedBy: "Scheduler-service",
				PublishedAt: timestamppb.Now(),
				Feature:     work.Feature,
				Data:        workData,
				Type:        "notification",
			}
			msgByte, err := proto.Marshal(kafkaMessage)
			if err != nil {
				log.Printf("weeklyWorkCommand Error: failed on create json message of ref2(%v)and owner(%v): %v ,", work.Ref2, work.Owner, err.Error())
				return
			}
			kafkaMsg := message.NewMessage(watermill.NewUUID(), msgByte)
			if err := kafkaPubliser.Publish(work.CallbackTopic, kafkaMsg); err != nil {
				log.Println("weeklyWorkCommand Error: failed on publish message: " + err.Error())
			}
		}
	}
}
