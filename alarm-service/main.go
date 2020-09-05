package main

import (
	"log"

	"github.com/Planxnx/message-processing-api/alarm-service/pkg/connection"

	scheduleRepository "github.com/Planxnx/message-processing-api/alarm-service/pkg/schedule/repository"
	scheduleUsecase "github.com/Planxnx/message-processing-api/alarm-service/pkg/schedule/usecase"
)

func main() {
	connection := connection.InitializeConnection()

	scheduleRepo := scheduleRepository.NewScheduleRepository(connection.MessageProcssingAPIDatabase.Collection("workSchedule"))
	scheduleUsecase := scheduleUsecase.NewScheduleUsecase(scheduleRepo)
	scheduleUsecase.CreateNewWorkSchedule()
	log.Println("Start Alarm-service...")
}
