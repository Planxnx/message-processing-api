package health

import (
	"context"
	"log"
	"time"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

var healthDataCaches = make(map[string]*HealthData)


type HealthUsercase struct {
	healthCollection    *qmgo.Collection
	healthLogCollection *qmgo.Collection
}

func New(mc *qmgo.Collection, hlmc *qmgo.Collection) *HealthUsercase {
	return &HealthUsercase{
		healthCollection:    mc,
		healthLogCollection: hlmc,
	}
}

func (pU *HealthUsercase) GetHealthByFeatureAndServiceName(ctx context.Context, feature string, serviceName string) (*HealthData, error) {
	healthData := &HealthData{}
	err := pU.healthCollection.Find(ctx, bson.M{
		"feature":     feature,
		"serviceName": serviceName,
	}).One(healthData)
	if err != nil {
		return nil, err
	}
	log.Println(healthData)

	return healthData, nil
}

func (pU *HealthUsercase) GetAllHealths(ctx context.Context) ([]*HealthData, error) {
	healthData := &[]*HealthData{}
	err := pU.healthCollection.Find(ctx, bson.M{}).All(healthData)
	if err != nil {
		return nil, err
	}
	log.Println(healthData)

	return *healthData, nil
}

func (*HealthUsercase) GetHealthMem(feature string) (*HealthData) {
	return healthDataCaches[feature]
}

func (pU *HealthUsercase) UpsertHealthData(ctx context.Context, healthData *HealthData) error {
	healthData.LastCheckedAt = time.Now()
	_, err := pU.healthCollection.Upsert(ctx,
		bson.M{
			"feature":     healthData.Feature,
			"serviceName": healthData.ServiceName,
		}, healthData)
	if err != nil {
		return err
	}
	go pU.createHealthDataLog(ctx, healthData)
	
	//save mem
	healthDataCaches[healthData.Feature] = healthData
	return nil
}

func (pU *HealthUsercase) createHealthDataLog(ctx context.Context, healthData *HealthData) error {
	_, err := pU.healthLogCollection.InsertOne(ctx, &HealthLog{
		Feature:     healthData.Feature,
		Description: healthData.Description,
		ExecuteMode: healthData.ExecuteMode,
		ServiceName: healthData.ServiceName,
		CheckedAt:   healthData.LastCheckedAt,
	})
	if err != nil {
		return err
	}
	return nil
}
