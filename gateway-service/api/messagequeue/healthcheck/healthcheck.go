package healthcheck

import (
	"context"
	"encoding/json"

	healthusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/health"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"

	"github.com/ThreeDotsLabs/watermill/message"
)

type HealthCheckHandler struct {
	healthUsercase *healthusecase.HealthUsercase
}

func New(h *healthusecase.HealthUsercase) *HealthCheckHandler {
	return &HealthCheckHandler{
		healthUsercase: h,
	}
}

func (h *HealthCheckHandler) HealthCheck(msg *message.Message) error {
	ctx := context.Background()
	resultMsg := &messageschema.HealthCheckMessageFormat{}
	json.Unmarshal(msg.Payload, resultMsg)

	err := h.healthUsercase.UpsertHealthData(ctx, &healthusecase.HealthData{
		Feature:     resultMsg.Feature,
		Description: resultMsg.Description,
		ExecuteMode: h.mapExcuteModeToString(resultMsg.ExecuteMode),
		ServiceName: resultMsg.ServiceName,
	})

	if err != nil {
		return err
	}

	return nil
}

func (*HealthCheckHandler) mapExcuteModeToString(execMode []messageschema.ExecuteMode) []string {
	execModeString := []string{}
	for _, v := range execMode {
		execModeString = append(execModeString, string(v))
	}
	return execModeString
}
