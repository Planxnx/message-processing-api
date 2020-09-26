package restful

import (
	"github.com/Planxnx/message-processing-api/botnoi-service/internal/api/restful/chitchat"
	"github.com/Planxnx/message-processing-api/botnoi-service/internal/botnoi"
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
