package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/Planxnx/message-processing-api/scheduler-service/internal/alarm"
	"github.com/Planxnx/message-processing-api/scheduler-service/internal/api/messagequeue"
	mqmessage "github.com/Planxnx/message-processing-api/scheduler-service/internal/api/messagequeue/message"
	"github.com/Planxnx/message-processing-api/scheduler-service/internal/message"
	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/connection"
	pkgcron "github.com/Planxnx/message-processing-api/scheduler-service/pkg/cron"
	scheduleRepository "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/repository"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/robfig/cron"
	"google.golang.org/protobuf/proto"
)

func init() {
	os.Setenv("TZ", "Asia/Bangkok")
}

var serviceName = "Scheduler-servic"

func main() {
	ctx := context.Background()

	connection := connection.InitializeConnection()

	scheduleRepo := scheduleRepository.NewScheduleRepository(connection.MessageProcssingAPIDatabase.Collection("workSchedule"))
	scheduleCron := pkgcron.NewScheduleUsecase(scheduleRepo, connection.KafkaPublisher)

	alarmService := alarm.New(scheduleRepo)
	messageUsecase := message.New(connection.KafkaPublisher)

	messageMQHandler := mqmessage.New(messageUsecase, alarmService)

	messageQueueRouterDependency := &messagequeue.RouterDependency{
		KafkaSubscriber: connection.KafkaSubscriber,
		MessageHandler:  messageMQHandler,
	}

	log.Printf("%s: Start cron job ...\n", serviceName)
	cron := scheduleCron.StartFetchSchedule()

	messagequeueRouter, err := messageQueueRouterDependency.InitialRouter()
	if err != nil {
		log.Fatalf("main Error: failed on create new messagequeue router: %v", err)
	}

	go healthCheck(serviceName, connection.KafkaPublisher)

	go func() {
		log.Printf("%s: Start messagequeue subscriber ...\n", serviceName)
		if err := messagequeueRouter.Run(ctx); err != nil {
			log.Fatalf("main Error: failed on start messagequeue subscruber: %v", err)
		}
	}()

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	<-killSignal

	cron.Stop()
	messagequeueRouter.Close()

	log.Println("Shutdown Scheduler-service...")
}

func healthCheck(serviceName string, kafkaPublisher *kafka.Publisher) {
	healthCheckCmd := func() {
		//Daily Notification HealthCheck
		go func() {
			dailynoti := &messageschema.HealthCheckMessage{
				Feature:     "Daily-Notification",
				Description: "แจ้งเตือนรายวัน ทุกวัน วันล่ะครั้ง",
				ExecuteMode: []messageschema.ExecuteMode{
					messageschema.ExecuteMode_Asynchronous,
				},
				ServiceName: serviceName,
			}
			dailynotiByte, err := proto.Marshal(dailynoti)
			if err != nil {
				log.Println("health check error: can't marshal dailynoti message")
			}
			if err := kafkaPublisher.Publish(messageschema.HealthCheckTopic, watermillmessage.NewMessage(watermill.NewShortUUID(), dailynotiByte)); err != nil {
				log.Printf("health check error: failed on publish dailynoti message: %v\n", err)
			}
		}()
	}

	//startup
	go healthCheckCmd()

	c := cron.New()
	c.AddFunc("@every 3m", healthCheckCmd)
	c.Start()
}
