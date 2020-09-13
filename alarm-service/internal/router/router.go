package router

import (
	"time"

	kafkaPkg "github.com/Planxnx/message-processing-api/alarm-service/pkg/connection/kafka"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"

	createAlarmDelivery "github.com/Planxnx/message-processing-api/alarm-service/internal/createAlarm/delivery"
)

var (
	logger = watermill.NewStdLogger(false, true)
)

func NewEventRouter(kafkaSubscriber *kafka.Subscriber, caDelivery *createAlarmDelivery.CreateAlarmDelivery) (*message.Router, error) {
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
	router.AddNoPublisherHandler("CommonCreateAlarmHandler", kafkaPkg.CommonMessage, kafkaSubscriber, caDelivery.CommonMessageHandler)
	return router, nil
}
