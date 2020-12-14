package mongodb

import (
	"context"

	"github.com/qiniu/qmgo"
)

func NewClient(ctx context.Context) (*qmgo.Client, error) {
	config := &qmgo.Config{
		Uri: "mongodb://admin:admin@localhost:27017",
	}
	return qmgo.NewClient(ctx, config)
}
