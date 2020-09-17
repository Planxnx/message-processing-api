package message

import (
	"encoding/json"
	"log"

	"github.com/Planxnx/message-processing-api/gateway-service/model"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func MessageHandle(ctx *fasthttp.RequestCtx) {
	reqBody := &model.MessageRequest{}
	json.Unmarshal(ctx.Request.Body(), reqBody)
	log.Printf("Received Message: %v", reqBody)

	resBody := &model.MessageResponse{
		Message: "Success",
		Data: model.MessageResponseData{
			MessageRef: uuid.New().String(),
		},
	}
	resJson, err := json.Marshal(resBody)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(resJson)
}
