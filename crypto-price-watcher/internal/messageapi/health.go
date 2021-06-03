package messageapi

import (
	"log"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill"
	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/protobuf/proto"
)

func (m *MessageClient) healthCheck(h *handler) {
	execMode := make([]messageschema.ExecuteMode, 0)
	if h.Config.AsynchronousMode {
		execMode = append(execMode, messageschema.ExecuteMode_Asynchronous)
	}
	if h.Config.SynchronousMode {
		execMode = append(execMode, messageschema.ExecuteMode_Synchronous)

	}

	healthByte, err := proto.Marshal(&messageschema.HealthCheckMessage{
		Feature:     h.Config.FeatureName,
		Description: h.Config.Description,
		ExecuteMode: execMode,
		ServiceName: h.Config.ServiceName,
	})
	if err != nil {
		log.Println("health check error: can't marshal message")
	}

	if err := m.KafkaPublisher.Publish(messageschema.HealthCheckTopic, watermillmessage.NewMessage(watermill.NewShortUUID(), healthByte)); err != nil {
		log.Printf("health check error: failed on publish message: %v\n", err)
	}
}
