package usecase

import (
	"context"
	"log"
	"time"

	scheduleConst "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/constant"
	scheduleModel "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/model"
	scheduleRepository "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/repository"
)

type CreateAlarmUsecase struct {
	ScheduleRepository *scheduleRepository.ScheduleRepository
}

func NewCreateAlarmUsecase(schRepo *scheduleRepository.ScheduleRepository) *CreateAlarmUsecase {
	return &CreateAlarmUsecase{
		ScheduleRepository: schRepo,
	}
}

func (ca *CreateAlarmUsecase) CreateDailyAlarm(ctx context.Context) error {
	tnow := time.Now()
	result, err := ca.ScheduleRepository.InsertWeeklySchedule(ctx, scheduleModel.WorkSchedule{
		Ref1:    "line_1",
		Ref2:    "message_1",
		Ref3:    "user_1",
		Owner:   "alarm_service",
		Message: "แจ้งเตือนแสตนด์อัพมีทติ้ง2",
		CallbackTopic: []string{
			"replyMessage",
			"boitnoiMessage",
		},
		Time: scheduleModel.WorkTime{
			Timestamp: tnow,
			Day:       10,
			WeekDay:   scheduleConst.WeekDay_FRIDAY,
			Hour:      02,
			Minute:    12,
			Second:    30,
		},
		Data: map[string]interface{}{
			"message": "ไปกินขี้นะมึงแ่พ",
			"score":   10,
		},
		Features: map[string]bool{
			"replyMessage": true,
			"dailyAlarm":   false,
		},
	})
	if err != nil {
		log.Printf("Error: failed on insert: %s", err.Error())
		return err
	}
	log.Printf("Insert successful: %v", result.InsertedID)
	return nil
}
