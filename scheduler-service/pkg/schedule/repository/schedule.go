package repository

import (
	"context"
	"time"

	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/schedule/constant"
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
		"type":        constant.ScheduleType_DAILY,
		"time.hour":   hh,
		"time.minute": mm,
	}).All(workSchedule)
	if err != nil {
		return nil, err
	}
	return workSchedule, nil
}

func (schRepo *ScheduleRepository) GetEveryHourSchedule(ctx context.Context) (*[]model.WorkSchedule, error) {
	tNow := time.Now()
	_, mm, _ := tNow.Clock()
	workSchedule := &[]model.WorkSchedule{}
	err := schRepo.WorkScheduleCollection.Find(ctx, bson.M{
		"type":        constant.ScheduleType_EVERYHOUR,
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
		"type": constant.ScheduleType_DAILY,
	}).All(workSchedule)
	if err != nil {
		return nil, err
	}
	return workSchedule, nil
}

func (schRepo *ScheduleRepository) GetAllEveryHourSchedule(ctx context.Context) (*[]model.WorkSchedule, error) {
	workSchedule := &[]model.WorkSchedule{}
	err := schRepo.WorkScheduleCollection.Find(ctx, bson.M{
		"type": constant.ScheduleType_EVERYHOUR,
	}).All(workSchedule)
	if err != nil {
		return nil, err
	}
	return workSchedule, nil
}

func (schRepo *ScheduleRepository) InsertDailySchedule(ctx context.Context, workSchedule model.WorkSchedule) (*qmgo.InsertOneResult, error) {
	workSchedule.Type = constant.ScheduleType_DAILY
	return schRepo.WorkScheduleCollection.InsertOne(ctx, workSchedule)
}

func (schRepo *ScheduleRepository) InsertEveryHourSchedule(ctx context.Context, workSchedule model.WorkSchedule) (*qmgo.InsertOneResult, error) {
	workSchedule.Type = constant.ScheduleType_EVERYHOUR
	return schRepo.WorkScheduleCollection.InsertOne(ctx, workSchedule)
}
