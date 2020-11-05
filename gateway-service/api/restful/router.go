package restful

import (
	"github.com/Planxnx/message-processing-api/gateway-service/api/restful/health"
	"github.com/Planxnx/message-processing-api/gateway-service/api/restful/message"
	"github.com/Planxnx/message-processing-api/gateway-service/api/restful/middleware"
	"github.com/Planxnx/message-processing-api/gateway-service/api/restful/provider"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
)

type RouterDependency struct {
	App            *fiber.App
	Middleware     *middleware.Middleware
	HealthHandler  *health.HealthHandler
	MessageHandler *message.MessageHandler
	ProvderHandler *provider.ProviderHandler
}

func (r *RouterDependency) InitialRouter() {

	r.App.Get("/health", r.HealthHandler.CheckHealth)

	v1 := r.App.Group("/v1")

	v1.Use(cors.New())
	v1.Use(logger.New())
	v1.Use(helmet.New())

	v1.Post("/provider/register", r.ProvderHandler.RegisterEndpoint)

	authRoute := v1.Group("/", r.Middleware.AuthenticationMiddleware)

	authRoute.Post("/", r.MessageHandler.MainEndpoint)
}
