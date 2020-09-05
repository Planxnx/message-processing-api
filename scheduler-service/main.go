package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/connection"
	scheduleRepository "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/repository"
	scheduleUsecase "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/usecase"
)

func main() {
	connection := connection.InitializeConnection()

	scheduleRepo := scheduleRepository.NewScheduleRepository(connection.MessageProcssingAPIDatabase.Collection("workSchedule"))
	scheduleUsecase := scheduleUsecase.NewScheduleUsecase(scheduleRepo)

	go scheduleUsecase.StartFetchSchedule()

	log.Println("Start Scheduler-service...")

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	<-killSignal
	log.Println("Shutdown Scheduler-service...")
}
