package restful

import (
	"github.com/Planxnx/message-processing-api/botnoi-service/internal/api/restful/chitchat"
	"github.com/buaazp/fasthttprouter"
)

func New() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.GET("/api/chitchat", chitchat.Handler)

	return router
}
