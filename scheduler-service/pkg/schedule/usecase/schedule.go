package schedule

import (
	"context"
	"log"

	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/repository"

	"github.com/robfig/cron/v3"
)

type ScheduleUsecase struct {
	ScheduleRepository *repository.ScheduleRepository
	Command            func()
	Cron               *cron.Cron
	CronSpec           string
}

func NewScheduleUsecase(schRepo *repository.ScheduleRepository) *ScheduleUsecase {
	return &ScheduleUsecase{
		ScheduleRepository: schRepo,
		Cron:               cron.New(),
		CronSpec:           "@every 10s",
	}
}

func (sch *ScheduleUsecase) StartFetchSchedule() {
	sch.Cron.AddFunc(sch.CronSpec, func() {
		log.Println("Fetching Schedule!")
		ctx := context.Background()
		workSch, _ := sch.ScheduleRepository.GetDailySchedule(ctx)
		log.Println(workSch)
	})
	sch.Cron.Start()
	log.Println("Start fetch schedule")
}
