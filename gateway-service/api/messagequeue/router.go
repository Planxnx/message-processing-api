package messagequeue

import (
	"time"

	healthcheckhandler "github.com/Planxnx/message-processing-api/gateway-service/api/messagequeue/healthcheck"
	messagehandler "github.com/Planxnx/message-processing-api/gateway-service/api/messagequeue/message"
	messageschema "github.com/Planxnx/message-processing-api/message-schema"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
)

var (
	logger = watermill.NewStdLogger(false, true)
)

type RouterDependency struct {
	KafkaSubscriber    *kafka.Subscriber
	MessageHandler     *messagehandler.MessageHandler
	HealthCheckHandler *healthcheckhandler.HealthCheckHandler
}

func (r *RouterDependency) InitialRouter() (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,
		middleware.Recoverer,
	)

	router.AddNoPublisherHandler("HealthCheckHandler", messageschema.HealthCheck, r.KafkaSubscriber, r.HealthCheckHandler.HealthCheck)
	router.AddNoPublisherHandler("ReplyMessageHandler", messageschema.ReplyMessage, r.KafkaSubscriber, r.MessageHandler.ReplyMessage)
	return router, nil
}
