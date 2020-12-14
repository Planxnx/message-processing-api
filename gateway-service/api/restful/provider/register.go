package provider

import (
	providerusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/provider"
	"github.com/Planxnx/message-processing-api/gateway-service/model"
	"github.com/gofiber/fiber/v2"
)

func (pH *ProviderHandler) RegisterEndpoint(c *fiber.Ctx) error {
	reqBody := &model.ProviderResgisterRequest{}
	c.BodyParser(reqBody)

	_, err := pH.ProviderUsecase.GetProviderByID(c.Context(), reqBody.ProviderID)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(&model.Response{
			Message: "id is already existed",
		})
	}

	result, err := pH.ProviderUsecase.CreateNewProvider(c.Context(), &providerusecase.ProviderData{
		ID:      reqBody.ProviderID,
		Name:    reqBody.ProviderName,
		Webhook: reqBody.Webhook,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&model.Response{
			Message: "internal server error :" + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&model.Response{
		Message: "Success",
		Data: &model.ProviderResgisterResponseData{
			ProviderID: result.ID,
			Token:      result.Token,
		},
	})
}
