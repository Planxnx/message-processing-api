package repository

import (
	"context"
	"log"
	"time"

	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/model"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ScheduleRepository struct {
	WorkScheduleCollection *qmgo.Collection
}

func NewScheduleRepository(workSch *qmgo.Collection) *ScheduleRepository {

	return &ScheduleRepository{
		WorkScheduleCollection: workSch,
	}
}

func (schRepo *ScheduleRepository) GetDailySchedule(ctx context.Context) (*[]model.WorkSchedule, error) {
	tNow := time.Now()
	hh, mm, _ := tNow.Clock()
	workSchedule := &[]model.WorkSchedule{}
	err := schRepo.WorkScheduleCollection.Find(ctx, bson.M{
		"type": "DAILY",
		"time": bson.M{
			"hour":   hh,
			"minute": mm,
		},
	}).All(workSchedule)
	if err != nil {
		return nil, err
	}
	log.Println(workSchedule)
	return workSchedule, nil
}
