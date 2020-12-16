module github.com/Planxnx/message-processing-api/scheduler-service

go 1.14

require (
	github.com/Planxnx/message-processing-api/external-caller-service v0.0.0-20201214223757-1e1d612e1f29
	github.com/Planxnx/message-processing-api/message-schema v0.0.0-20201214223757-1e1d612e1f29
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/ThreeDotsLabs/watermill-kafka/v2 v2.2.0
	github.com/pkg/errors v0.9.1
	github.com/qiniu/qmgo v0.7.0
	github.com/robfig/cron v1.2.0
	github.com/robfig/cron/v3 v3.0.1
	go.mongodb.org/mongo-driver v1.4.1
	google.golang.org/protobuf v1.25.0
)
