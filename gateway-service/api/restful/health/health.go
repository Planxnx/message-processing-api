package health

import "github.com/gofiber/fiber/v2"

type HealthHandler struct{}

func New() *HealthHandler {
	return &HealthHandler{}
}

func (HealthHandler) CheckHealth(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Health OK")
}
