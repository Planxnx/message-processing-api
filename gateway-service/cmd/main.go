package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Planxnx/message-processing-api/gateway-service/api/restful"
	"github.com/Planxnx/message-processing-api/gateway-service/api/restful/health"
	"github.com/Planxnx/message-processing-api/gateway-service/api/restful/message"

	mqmessage "github.com/Planxnx/message-processing-api/gateway-service/api/messagequeue/message"

	messageusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/message"

	"github.com/Planxnx/message-processing-api/gateway-service/model"
	"github.com/gofiber/fiber/v2"

	"github.com/Planxnx/message-processing-api/gateway-service/api/messagequeue"
	kafapkg "github.com/Planxnx/message-processing-api/gateway-service/pkg/kafka"
)

var (
	port string = "8080"
	wg   sync.WaitGroup
)

func init() {
	os.Setenv("TZ", "Asia/Bangkok")
}

func main() {
	ctx := context.Background()

	app := fiber.New(fiber.Config{
		ErrorHandler: defaultErrorHandler,
	})

	//Initial Dependency
	kafkaSubscriber, err := kafapkg.NewSubscriber()
	if err != nil {
		log.Fatalf("main Error: failed on create kafka subscriber: %v", err)
	}
	kafkaNewPublisher, err := kafapkg.NewPubliser()
	if err != nil {
		log.Fatalf("main Error: failed on create kafka publisher: %v", err)
	}
	messageUsecase := messageusecase.New(kafkaNewPublisher)

	//Initial MessageQueue Dependency
	messageMQHandler := mqmessage.New(messageUsecase)
	messageQueueouterDependency := &messagequeue.RouterDependency{
		KafkaSubscriber: kafkaSubscriber,
		MessageHandler:  messageMQHandler,
	}
	messagequeueRouter, err := messageQueueouterDependency.InitialRouter()
	if err != nil {
		log.Fatalf("main Error: failed on create new messagequeue router: %v", err)
	}

	//Initial Restful Dependency
	healthHandler := health.New()
	messageHandler := message.New(messageUsecase)
	routerDependency := &restful.RouterDependency{
		HealthHandler:  healthHandler,
		App:            app,
		MessageHandler: messageHandler,
	}
	routerDependency.InitialRouter()

	//Start services
	go func() {
		log.Println("Start messagequeue subscriber ...")
		if err := messagequeueRouter.Run(ctx); err != nil {
			log.Fatalf("main Error: failed on start messagequeue subscruber: %v", err)
		}
	}()
	go func() {
		log.Println("Start restful server on :" + port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("main Error: failed on start restful server: %v", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}
}

func defaultErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(model.Response{
		Message: err.Error(),
	})
}
