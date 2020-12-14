module github.com/Planxnx/message-processing-api/gateway-service

go 1.15

replace github.com/Planxnx/message-processing-api/message-schema => ../message-schema

require (
	github.com/Planxnx/message-processing-api/message-schema v0.0.0-20200914043633-32f835e1f4da
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/ThreeDotsLabs/watermill-kafka/v2 v2.2.0
	github.com/gofiber/fiber/v2 v2.0.1
	github.com/gofiber/helmet/v2 v2.0.0
	github.com/google/uuid v1.1.2 // indirect
	github.com/qiniu/qmgo v0.7.0
	github.com/valyala/fasthttp v1.16.0
	go.mongodb.org/mongo-driver v1.4.1
	google.golang.org/protobuf v1.25.0
)
