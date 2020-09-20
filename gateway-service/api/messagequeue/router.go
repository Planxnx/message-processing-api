package messagequeue

import (
	"encoding/json"
	"log"
	"time"

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

func NewMessageQueueRouter(kafkaSubscriber *kafka.Subscriber) (*message.Router, error) {
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

	router.AddNoPublisherHandler("ReplyMessageHandler", messageSchema.ReplyMessage, kafkaSubscriber, func(msg *message.Message) error {
		resultMsg := &messageSchema.DefaultMessageFormat{}
		json.Unmarshal(msg.Payload, resultMsg)

		log.Printf("Received Notification Event:\n  ---Unmarshal Message: %v \n  ---Raw Message: %v \n", resultMsg, string(msg.Payload))

		return nil
	})
	return router, nil
}
