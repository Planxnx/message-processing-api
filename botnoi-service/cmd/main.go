package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Planxnx/message-processing-api/botnoi-service/config"
	"github.com/Planxnx/message-processing-api/botnoi-service/internal/api/restful"
	"github.com/valyala/fasthttp"
)

func main() {
	configs := config.InitialConfig()
	restfulAPI := restful.New()
	log.Println("start server on :", configs.Restful.Port)
	if err := fasthttp.ListenAndServe(":"+strconv.Itoa(configs.Restful.Port), restfulAPI.Handler); err != http.ErrServerClosed {
		log.Fatalf("main Error: failed on start server: %v", err)
	}
}
