package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/connection"
	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule"
)

func main() {
	connection := connection.InitializeConnection()
	scheduler := schedule.InitializeSchedule(connection)
	go scheduler.StartSchedule()
	log.Println("Start Scheduler-service...")

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	<-killSignal
	log.Println("Shutdown Scheduler-service...")
}
