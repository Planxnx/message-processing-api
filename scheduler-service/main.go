package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/connection"
	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/cron"
	scheduleRepository "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/repository"
)

func main() {
	connection := connection.InitializeConnection()

	scheduleRepo := scheduleRepository.NewScheduleRepository(connection.MessageProcssingAPIDatabase.Collection("workSchedule"))
	scheduleCron := cron.NewScheduleUsecase(scheduleRepo)

	go scheduleCron.StartFetchSchedule()

	log.Println("Start Scheduler-service...")

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	<-killSignal
	log.Println("Shutdown Scheduler-service...")
}
