package main

import (
	"context"
	"log"

	"github.com/Planxnx/message-processing-api/external-caller-service/config"
	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/protobuf/proto"

	"github.com/Planxnx/message-processing-api/external-caller-service/internal/api/messagequeue"
	mqmessage "github.com/Planxnx/message-processing-api/external-caller-service/internal/api/messagequeue/message"
	"github.com/Planxnx/message-processing-api/external-caller-service/internal/botnoi"
	"github.com/Planxnx/message-processing-api/external-caller-service/internal/lottery"
	"github.com/Planxnx/message-processing-api/external-caller-service/internal/message"

	kafapkg "github.com/Planxnx/message-processing-api/external-caller-service/pkg/kafka"
	"github.com/robfig/cron/v3"
)

func main() {
	ctx := context.Background()

	configs := config.InitialConfig()

	//Initial Dependency
	kafkaSubscriber, err := kafapkg.NewSubscriber()
	if err != nil {
		log.Fatalf("main Error: failed on create kafka subscriber: %v", err)
	}
	kafkaNewPublisher, err := kafapkg.NewPubliser()
	if err != nil {
		log.Fatalf("main Error: failed on create kafka publisher: %v", err)
	}

	//Initial Usecase
	botnoiUsecase := botnoi.New(configs.Botnoi.Address, configs.Botnoi.Token)
	messageUsecase := message.NewUsecase(kafkaNewPublisher)
	lottoUsecase := lottery.New(configs.Lottery.Address)

	messageMQHandler := mqmessage.New(messageUsecase, botnoiUsecase, lottoUsecase)

	messageQueueRouterDependency := &messagequeue.RouterDependency{
		KafkaSubscriber: kafkaSubscriber,
		MessageHandler:  messageMQHandler,
	}
	messagequeueRouter, err := messageQueueRouterDependency.InitialRouter()
	if err != nil {
		log.Fatalf("main Error: failed on create new messagequeue router: %v", err)
	}

	go healthCheck(configs.ServiceName, kafkaNewPublisher)

	log.Printf("%s: Start messagequeue subscriber ...\n", configs.ServiceName)
	if err := messagequeueRouter.Run(ctx); err != nil {
		log.Fatalf("main Error: failed on start messagequeue subscruber: %v", err)
	}
}

func healthCheck(serviceName string, kafkaPublisher *kafka.Publisher) {
	healthCheckCmd := func() {
		//Chitchat HealthCheck
		go func() {
			chitchat := &messageschema.HealthCheckMessage{
				Feature:     "Chitchat",
				Description: "แชทบอทคุยเล่นขำขัน",
				ExecuteMode: []messageschema.ExecuteMode{
					messageschema.ExecuteMode_Asynchronous,
					messageschema.ExecuteMode_Synchronous,
				},
				ServiceName: serviceName,
			}
			chitchatByte, err := proto.Marshal(chitchat)
			if err != nil {
				log.Println("health check error: can't marshal chitchat message")
			}
			if err := kafkaPublisher.Publish(messageschema.HealthCheckTopic, watermillmessage.NewMessage(watermill.NewShortUUID(), chitchatByte)); err != nil {
				log.Printf("health check error: failed on publish chitchat message: %v\n", err)
			}
		}()

		//CheclLatestLottery HealthCheck
		go func() {
			lotto := &messageschema.HealthCheckMessage{
				Feature:     "Check-Latest-Lottery",
				Description: "ตรวจผลสลากกินแบ่งรัฐบาล งวดล่าสุด",
				ExecuteMode: []messageschema.ExecuteMode{
					messageschema.ExecuteMode_Synchronous,
				},
				ServiceName: serviceName,
			}
			chitchatByte, err := proto.Marshal(lotto)
			if err != nil {
				log.Println("health check error: can't marshal lotto message")
			}
			if err := kafkaPublisher.Publish(messageschema.HealthCheckTopic, watermillmessage.NewMessage(watermill.NewShortUUID(), chitchatByte)); err != nil {
				log.Printf("health check error: failed on publish lotto message: %v\n", err)
			}
		}()
	}

	//startup
	go healthCheckCmd()

	c := cron.New()
	c.AddFunc("@every 3m", healthCheckCmd)
	c.Start()
}
