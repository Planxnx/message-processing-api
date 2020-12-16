package alarm

import (
	scheduleRepository "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/repository"
)

type Service struct {
	scheduleRepository *scheduleRepository.ScheduleRepository
}

func New(sR *scheduleRepository.ScheduleRepository) *Service {
	return &Service{
		scheduleRepository: sR,
	}
}
