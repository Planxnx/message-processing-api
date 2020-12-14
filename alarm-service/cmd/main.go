package main

import (
	"context"
	"log"

	createAlarmDelivery "github.com/Planxnx/message-processing-api/alarm-service/internal/createAlarm/delivery"
	createAlarmUsecase "github.com/Planxnx/message-processing-api/alarm-service/internal/createAlarm/usecase"
	eventRouter "github.com/Planxnx/message-processing-api/alarm-service/internal/router"
	"github.com/Planxnx/message-processing-api/alarm-service/pkg/connection"
	scheduleRepository "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/repository"
)

func main() {
	ctx := context.Background()

	connection := connection.InitializeConnection()
	scheduleRepo := scheduleRepository.NewScheduleRepository(connection.MessageProcssingAPIDatabase.Collection("workSchedule"))
	createAlarmUseC := createAlarmUsecase.NewCreateAlarmUsecase(scheduleRepo)
	createAlarmDeli := createAlarmDelivery.NewCreateAlarmDelivery(createAlarmUseC)

	router, err := eventRouter.NewEventRouter(connection.KafkaSubscriber, createAlarmDeli)
	if err != nil {
		panic(err)
	}

	log.Println("Start Alarm-service...")
	if err := router.Run(ctx); err != nil {
		panic(err)
	}

}
