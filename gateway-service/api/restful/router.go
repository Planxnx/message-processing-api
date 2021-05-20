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

	r.App.Use(cors.New())
	r.App.Use(logger.New())
	r.App.Use(helmet.New())

	r.App.Get("/health", r.HealthHandler.CheckHealth)
	r.App.Get("/health/:service/:feature", r.HealthHandler.CheckHealthByFeatureAndService)

	v1 := r.App.Group("/v1")

	v1.Post("/provider/register", r.ProvderHandler.RegisterEndpoint)

	authRoute := v1.Group("/", r.Middleware.AuthenticationMiddleware)

	authRoute.Post("/", r.MessageHandler.MainEndpoint)

	authRoute.Post("/sync", r.MessageHandler.SynchronousEndpoint)
}
