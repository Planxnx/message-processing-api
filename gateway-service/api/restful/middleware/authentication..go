package middleware

import (
	"github.com/gofiber/fiber/v2"

	providerusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/provider"
	"github.com/Planxnx/message-processing-api/gateway-service/model"
)

type Middleware struct {
	ProviderUsecase *providerusecase.ProviderUsercase
}

func New(pU *providerusecase.ProviderUsercase) *Middleware {
	return &Middleware{
		ProviderUsecase: pU,
	}
}

func (m *Middleware) AuthenticationMiddleware(c *fiber.Ctx) error {
	providerID := c.Get("providerID")
	secret := c.Get("secret")

	providerData, err := m.ProviderUsecase.GetProviderByID(c.Context(), providerID)
	if err != nil || providerData.Secret != secret {
		return c.Status(fiber.StatusUnauthorized).JSON(&model.Response{
			Message: "provider not found or wrong secret",
		})
	}

	return c.Next()
}
