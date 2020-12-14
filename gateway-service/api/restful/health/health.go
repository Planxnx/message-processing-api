package health

import (
	"log"
	"time"

	healthusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/health"
	"github.com/Planxnx/message-processing-api/gateway-service/model"
	"github.com/gofiber/fiber/v2"
	"github.com/qiniu/qmgo"
)

type HealthHandler struct {
	healthUsercase *healthusecase.HealthUsercase
}

func New(healthUsercase *healthusecase.HealthUsercase) *HealthHandler {
	return &HealthHandler{
		healthUsercase: healthUsercase,
	}
}

func (h *HealthHandler) CheckHealth(c *fiber.Ctx) error {
	healths, err := h.healthUsercase.GetAllHealths(c.Context())
	if err != nil {
		log.Printf("CheckHealth Error: failed on get all healths: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	responseData := []*model.HealthCheckResponseData{}

	for _, health := range healths {
		healcheckdata := &model.HealthCheckResponseData{
			Feature:     health.Feature,
			Description: health.Description,
			ExecuteMode: health.ExecuteMode,
			ServiceName: health.ServiceName,
			Status:      true,
		}

		timeDiff := time.Now().Sub(health.LastCheckedAt)
		if health.LastCheckedAt.Before(time.Now().Add(-5 * time.Minute)) {
		}
		if timeDiff > time.Duration(5*time.Minute) {
			healcheckdata.Status = false
			healcheckdata.LastOnline = &health.LastCheckedAt
		}
		responseData = append(responseData, healcheckdata)
	}

	return c.Status(fiber.StatusOK).JSON(&model.Response{
		Message: "Success",
		Data:    responseData,
	})
}

func (h *HealthHandler) CheckHealthByFeatureAndService(c *fiber.Ctx) error {
	service := c.Params("service")
	feature := c.Params("feature")
	health, err := h.healthUsercase.GetHealthByFeatureAndServiceName(c.Context(), feature, service)
	if err != nil {
		if qmgo.IsErrNoDocuments(err) {
			return c.Status(fiber.StatusBadRequest).JSON(&model.Response{
				Message: "feature not found",
			})
		}
		log.Printf("CheckHealthByFeatureAndService Error: failed on get health: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	healcheckdata := &model.HealthCheckResponseData{
		Feature:     health.Feature,
		Description: health.Description,
		ExecuteMode: health.ExecuteMode,
		ServiceName: health.ServiceName,
		Status:      true,
	}

	timeDiff := time.Now().Sub(health.LastCheckedAt)
	if health.LastCheckedAt.Before(time.Now().Add(-5 * time.Minute)) {
	}
	if timeDiff > time.Duration(5*time.Minute) {
		healcheckdata.Status = false
		healcheckdata.LastOnline = &health.LastCheckedAt
	}

	return c.Status(fiber.StatusOK).JSON(&model.Response{
		Message: "Success",
		Data:    healcheckdata,
	})
}
