package connection

import (
	"context"
	"fmt"
	"log"

	kafaPkg "github.com/Planxnx/message-processing-api/scheduler-service/pkg/connection/kafka"
	"github.com/Planxnx/message-processing-api/scheduler-service/pkg/connection/mongodb"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/qiniu/qmgo"
)

type Connections struct {
	MongoDBClient               *qmgo.Client
	MessageProcssingAPIDatabase *qmgo.Database
	KafkaPublisher              *kafka.Publisher
	KafkaSubscriber             *kafka.Subscriber
}

func InitializeConnection() *Connections {
	log.Println("Initialize Conections")
	ctx := context.Background()

	mongodbclient, err := mongodb.NewClient(ctx)
	if err != nil {
		panic(fmt.Errorf("InitializeConnection Error: failed on create mongodb connection: %v", err.Error()))
	}

	kafkaPublisher, err := kafaPkg.NewPubliser()
	if err != nil {
		panic(fmt.Errorf("InitializeConnection Error: failed on create kafka connection: %v", err.Error()))
	}

	kafkaSubscriber, err := kafaPkg.NewSubscriber()
	if err != nil {
		panic(fmt.Errorf("InitializeConnection Error: failed on create kafka subscriber connection: %v", err.Error()))
	}

	return &Connections{
		MongoDBClient:               mongodbclient,
		MessageProcssingAPIDatabase: mongodbclient.Database("message-processing-api"),
		KafkaPublisher:              kafkaPublisher,
		KafkaSubscriber:             kafkaSubscriber,
	}
}
