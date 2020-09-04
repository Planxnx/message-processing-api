package connection

import (
	"log"
	"time"
)

type Mock struct {
}

func InitializeMock() *Mock {
	log.Println("Initialize Mock")

	return &Mock{}
}

func (Mock) GetNowSchedule() (WorkSchedule, error) {
	newTime, _ := time.Parse("2006-01-02", "2020-03-17")
	return WorkSchedule{
		Owner:   "Planx",
		Message: "ทดสอบจ้า",
		Time:    newTime,
	}, nil
}
