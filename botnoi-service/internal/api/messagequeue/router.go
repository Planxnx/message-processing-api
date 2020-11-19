package messagequeue

import (
	"time"

	messagehandler "github.com/Planxnx/message-processing-api/botnoi-service/internal/api/messagequeue/message"
	messageSchema "github.com/Planxnx/message-processing-api/message-schema"
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
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,
		middleware.Recoverer,
	)

	router.AddNoPublisherHandler("ReplyMessageHandler", messageSchema.CommonMessage, r.KafkaSubscriber, r.MessageHandler.ChitchatHandler)
	return router, nil
}
