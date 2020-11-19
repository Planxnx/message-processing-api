package main

import (
	"context"
	"log"

	"github.com/Planxnx/message-processing-api/botnoi-service/config"

	"github.com/Planxnx/message-processing-api/botnoi-service/internal/api/messagequeue"
	mqmessage "github.com/Planxnx/message-processing-api/botnoi-service/internal/api/messagequeue/message"
	"github.com/Planxnx/message-processing-api/botnoi-service/internal/botnoi"
	"github.com/Planxnx/message-processing-api/botnoi-service/internal/message"

	kafapkg "github.com/Planxnx/message-processing-api/botnoi-service/pkg/kafka"
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

	messageMQHandler := mqmessage.New(messageUsecase, botnoiUsecase)

	messageQueueRouterDependency := &messagequeue.RouterDependency{
		KafkaSubscriber: kafkaSubscriber,
		MessageHandler:  messageMQHandler,
	}
	messagequeueRouter, err := messageQueueRouterDependency.InitialRouter()
	if err != nil {
		log.Fatalf("main Error: failed on create new messagequeue router: %v", err)
	}

	log.Println("Start messagequeue subscriber ...")
	if err := messagequeueRouter.Run(ctx); err != nil {
		log.Fatalf("main Error: failed on start messagequeue subscruber: %v", err)
	}
}
