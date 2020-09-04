package connection

import (
	"context"
	"log"

	"github.com/qiniu/qmgo"
)

type Connections struct {
	MockCon       *Mock
	MongoDBClient *qmgo.Client
}

func InitializeConnection() *Connections {
	log.Println("Initialize Conections")
	ctx := context.Background()

	mongodbclient, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://admin:admin@localhost:27017"})
	if err != nil {
		panic(err)
	}

	return &Connections{
		MockCon:       InitializeMock(),
		MongoDBClient: mongodbclient,
	}
}
