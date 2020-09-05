package schedule

import (
	"context"
	"log"

	"github.com/Planxnx/message-processing-api/alarm-service/pkg/schedule/model"
	"github.com/Planxnx/message-processing-api/alarm-service/pkg/schedule/repository"
)

type CronCommand struct {
	CommandFunction func()
	TimeSpec        string
}

type ScheduleUsecase struct {
	ScheduleRepository *repository.ScheduleRepository
}

func NewScheduleUsecase(schRepo *repository.ScheduleRepository) *ScheduleUsecase {
	return &ScheduleUsecase{
		ScheduleRepository: schRepo,
	}
}

func (sch *ScheduleUsecase) CreateNewWorkSchedule() {
	ctx := context.Background()
	_, err := sch.ScheduleRepository.InsertDailySchedule(ctx, model.WorkSchedule{
		RefID:   "001",
		Owner:   "Planx",
		Topic:   "AlarmMessage",
		Message: "สวัสดีจ้า",
		Time: model.WorkTime{
			Date:   "2020-03-17",
			Hour:   21,
			Minute: 05,
			Second: 17,
		},
	})
	if err != nil {
		log.Println("ERR: " + err.Error())
		return
	}
	log.Println("Successful")
}
