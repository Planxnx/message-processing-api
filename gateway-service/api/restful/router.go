package restful

import (
	"github.com/Planxnx/message-processing-api/gateway-service/api/restful/health"
	"github.com/Planxnx/message-processing-api/gateway-service/api/restful/message"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
)

type RouterDependency struct {
	App            *fiber.App
	HealthHandler  *health.HealthHandler
	MessageHandler *message.MessageHandler
}

func (r *RouterDependency) InitialRouter() {

	r.App.Get("/health", r.HealthHandler.CheckHealth)

	v1 := r.App.Group("/v1")

	v1.Use(cors.New())
	v1.Use(logger.New())
	v1.Use(helmet.New())

	v1.Post("/", r.MessageHandler.MainEndpoint)
}
