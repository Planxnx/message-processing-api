package chitchat

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"

	"github.com/Planxnx/message-processing-api/external-caller-service/internal/botnoi"
)

type chitchatHandler struct {
	BotnoiService *botnoi.BotnoiService
}

type chitChatResponseBody struct {
	Message string `json:"message"`
}

type chitChatRequestBody struct {
	Message string `json:"message"`
}

func NewHandler(b *botnoi.BotnoiService) *chitchatHandler {
	return &chitchatHandler{
		BotnoiService: b,
	}
}

func (c *chitchatHandler) Handler(ctx *fasthttp.RequestCtx) {
	reqBody := &chitChatRequestBody{}
	json.Unmarshal(ctx.Request.Body(), reqBody)
	log.Println(reqBody)
	replyMessage, err := c.BotnoiService.ChitChatMessage(reqBody.Message)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	resp, _ := json.Marshal(&chitChatResponseBody{
		Message: replyMessage,
	})
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(resp)
	return
}
