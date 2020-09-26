package chitchat

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type chitchatHandler struct{}

func Handler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "CHITCHAT\n")
}
