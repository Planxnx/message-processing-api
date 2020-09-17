package api

import (
	"github.com/Planxnx/message-processing-api/gateway-service/api/health"
	"github.com/Planxnx/message-processing-api/gateway-service/api/message"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type RouterDependency struct {
}

func (RouterDependency) GetFastHTTPRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.GET("/health", fasthttpadaptor.NewFastHTTPHandlerFunc(health.HealthHandle))
	router.POST("/", message.MessageHandle)

	return router
}
