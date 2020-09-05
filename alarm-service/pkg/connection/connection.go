package connection

import (
	"context"
	"log"

	"github.com/qiniu/qmgo"
)

type Connections struct {
	MongoDBClient               *qmgo.Client
	MessageProcssingAPIDatabase *qmgo.Database
}

func InitializeConnection() *Connections {
	log.Println("Initialize Conections")
	ctx := context.Background()

	mongodbclient, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://admin:admin@localhost:27017"})
	if err != nil {
		panic(err)
	}

	return &Connections{
		MongoDBClient:               mongodbclient,
		MessageProcssingAPIDatabase: mongodbclient.Database("message-processing-api"),
	}
}
