package messagequeue

import (
	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	messagehandler "github.com/Planxnx/message-processing-api/scheduler-service/internal/api/messagequeue/message"
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
	KafkaSubscriber *kafka.Subscriber
	MessageHandler  *messagehandler.MessageHandler
}

func (r *RouterDependency) InitialRouter() (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Recoverer,
	)

	router.AddNoPublisherHandler("RegisterDailyNotificationHandler", messageschema.CommonMessageTopic, r.KafkaSubscriber, r.MessageHandler.RegisterDailyNotificationHandler)

	return router, nil
}
