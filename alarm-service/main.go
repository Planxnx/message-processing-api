package main

import (
	"context"
	"log"

	"github.com/Planxnx/message-processing-api/alarm-service/pkg/connection"

	scheduleModel "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/model"
	scheduleRepository "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/repository"
)

func main() {
	ctx := context.Background()
	connection := connection.InitializeConnection()

	scheduleRepo := scheduleRepository.NewScheduleRepository(connection.MessageProcssingAPIDatabase.Collection("workSchedule"))
	result, err := scheduleRepo.InsertEveryHourSchedule(ctx, scheduleModel.WorkSchedule{
		RefID:         "Test0000x",
		Owner:         "Planx",
		Message:       "แจ้งเตือนแสตนด์อัพมีทติ้ง",
		CallbackTopic: "AlarmMessage",
		Time: scheduleModel.WorkTime{
			Date:   "2020-03-17",
			Hour:   10,
			Minute: 30,
			Second: 30,
		},
	})
	if err != nil {
		log.Printf("Error: failed on insert: %s", err.Error())
		return
	}
	log.Printf("Insert successful: %v", result.InsertedID)
	log.Println("Start Alarm-service...")
}
