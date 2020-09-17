package main

import (
	"log"
	"os"

	"github.com/Planxnx/message-processing-api/gateway-service/api"
	"github.com/Planxnx/message-processing-api/gateway-service/api/health"
	"github.com/Planxnx/message-processing-api/gateway-service/api/message"

	messageusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/message"

	"github.com/Planxnx/message-processing-api/gateway-service/model"
	"github.com/gofiber/fiber/v2"
)

var (
	port string = "8080"
)

func init() {
	os.Setenv("TZ", "Asia/Bangkok")
}

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: defaultErrorHandler,
	})

	messageUsecase := messageusecase.New()

	healthHandler := health.New()
	messageHandler := message.New(messageUsecase)

	routerDependency := &api.RouterDependency{
		HealthHandler:  healthHandler,
		App:            app,
		MessageHandler: messageHandler,
	}
	routerDependency.InitialRouter()

	log.Println("start server on :" + port)
	log.Fatal(app.Listen(":" + port))
}

func defaultErrorHandler(ctx *fiber.Ctx, err error) error {
	return ctx.JSON(model.MessageResponse{
		Message: err.Error(),
	})
}
