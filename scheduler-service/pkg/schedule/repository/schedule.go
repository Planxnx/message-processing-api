package repository

import (
	"context"
	"log"
	"time"

	scheduleconstant "github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/constant"
	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/model"
	"github.com/pkg/errors"
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
		"type":        scheduleconstant.ScheduleType_DAILY,
		"time.hour":   hh,
		"time.minute": mm,
	}).All(workSchedule)
	if err != nil {
		return nil, err
	}
	return workSchedule, nil
}

func (schRepo *ScheduleRepository) GetWeeklySchedule(ctx context.Context) (*[]model.WorkSchedule, error) {
	tNow := time.Now()
	hh, mm, _ := tNow.Clock()
	workSchedule := &[]model.WorkSchedule{}
	log.Printf("%v %v: %v:%v", scheduleconstant.ScheduleType_WEEKLY, tNow.Weekday().String(), hh, mm)
	err := schRepo.WorkScheduleCollection.Find(ctx, bson.M{
		"type":         scheduleconstant.ScheduleType_WEEKLY,
		"time.weekDay": tNow.Weekday().String(),
		"time.hour":    hh,
		"time.minute":  mm,
	}).All(workSchedule)
	if err != nil {
		return nil, err
	}
	return workSchedule, nil
}

func (schRepo *ScheduleRepository) GetHOURLYSchedule(ctx context.Context) (*[]model.WorkSchedule, error) {
	tNow := time.Now()
	_, mm, _ := tNow.Clock()
	workSchedule := &[]model.WorkSchedule{}
	err := schRepo.WorkScheduleCollection.Find(ctx, bson.M{
		"type":        scheduleconstant.ScheduleType_HOURLY,
		"time.minute": mm,
	}).All(workSchedule)
	if err != nil {
		return nil, err
	}
	return workSchedule, nil
}

func (schRepo *ScheduleRepository) GetAllSchedule(ctx context.Context) (*[]model.WorkSchedule, error) {
	workSchedule := &[]model.WorkSchedule{}
	err := schRepo.WorkScheduleCollection.Find(ctx, bson.M{}).All(workSchedule)
	if err != nil {
		return nil, err
	}
	return workSchedule, nil
}

func (schRepo *ScheduleRepository) GetAllDailySchedule(ctx context.Context) (*[]model.WorkSchedule, error) {
	workSchedule := &[]model.WorkSchedule{}
	err := schRepo.WorkScheduleCollection.Find(ctx, bson.M{
		"type": scheduleconstant.ScheduleType_DAILY,
	}).All(workSchedule)
	if err != nil {
		return nil, err
	}
	return workSchedule, nil
}

func (schRepo *ScheduleRepository) GetAllHOURLYSchedule(ctx context.Context) (*[]model.WorkSchedule, error) {
	workSchedule := &[]model.WorkSchedule{}
	err := schRepo.WorkScheduleCollection.Find(ctx, bson.M{
		"type": scheduleconstant.ScheduleType_HOURLY,
	}).All(workSchedule)
	if err != nil {
		return nil, err
	}
	return workSchedule, nil
}

func (schRepo *ScheduleRepository) InsertiDailySchedule(ctx context.Context, workSchedule model.WorkSchedule) (*qmgo.InsertOneResult, error) {
	workSchedule.Type = scheduleconstant.ScheduleType_DAILY
	if workSchedule.Data == nil {
		workSchedule.Data = map[string]interface{}{}
	}
	workSchedule.CreateAt = time.Now()
	return schRepo.WorkScheduleCollection.InsertOne(ctx, workSchedule)
}

func (schRepo *ScheduleRepository) InsertHOURLYSchedule(ctx context.Context, workSchedule model.WorkSchedule) (*qmgo.InsertOneResult, error) {
	workSchedule.Type = scheduleconstant.ScheduleType_HOURLY
	if workSchedule.Data == nil {
		workSchedule.Data = map[string]interface{}{}
	}
	workSchedule.CreateAt = time.Now()
	return schRepo.WorkScheduleCollection.InsertOne(ctx, workSchedule)
}

func (schRepo *ScheduleRepository) InsertWeeklySchedule(ctx context.Context, workSchedule model.WorkSchedule) (*qmgo.InsertOneResult, error) {
	if workSchedule.Time.WeekDay == "" {
		return nil, errors.Errorf("WeekDay is invalid")
	}
	workSchedule.Type = scheduleconstant.ScheduleType_WEEKLY
	if workSchedule.Data == nil {
		workSchedule.Data = map[string]interface{}{}
	}
	if workSchedule.Data == nil {
		workSchedule.Data = map[string]interface{}{}
	}
	workSchedule.CreateAt = time.Now()
	return schRepo.WorkScheduleCollection.InsertOne(ctx, workSchedule)
}

func (schRepo *ScheduleRepository) GetWorkScheduleByReference(ctx context.Context, refType string, refId string, scheduleType string) (*[]model.WorkSchedule, error) {
	workSchedule := &[]model.WorkSchedule{}
	err := schRepo.WorkScheduleCollection.Find(ctx, bson.M{
		refType: refId,
		"type":  scheduleType,
	}).All(workSchedule)
	if err != nil {
		return nil, err
	}
	return workSchedule, nil
}

func midnightTimeConvert(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
