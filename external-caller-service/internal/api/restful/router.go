package restful

import (
	"github.com/Planxnx/message-processing-api/external-caller-service/internal/api/restful/chitchat"
	"github.com/Planxnx/message-processing-api/external-caller-service/internal/botnoi"
	"github.com/buaazp/fasthttprouter"
)

type RouterParamas struct {
	BotnoiService *botnoi.BotnoiService
}

func New(r *RouterParamas) *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.GET("/api/chitchat", chitchat.NewHandler(r.BotnoiService).Handler)

	return router
}
