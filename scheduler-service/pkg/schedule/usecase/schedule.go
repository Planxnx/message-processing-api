package schedule

import (
	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/repository"
)

type ScheduleUsecase struct {
	ScheduleRepository *repository.ScheduleRepository
}

func NewScheduleUsecase(schRepo *repository.ScheduleRepository) *ScheduleUsecase {
	return &ScheduleUsecase{
		ScheduleRepository: schRepo,
	}
}
