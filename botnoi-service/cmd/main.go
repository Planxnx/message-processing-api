package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Planxnx/message-processing-api/botnoi-service/config"
	"github.com/Planxnx/message-processing-api/botnoi-service/internal/api/restful"
	"github.com/valyala/fasthttp"

	"github.com/Planxnx/message-processing-api/botnoi-service/internal/botnoi"
)

func main() {
	configs := config.InitialConfig()

	botnoiService := botnoi.New(configs.Botnoi.Address, configs.Botnoi.Token)

	restfulAPI := restful.New(&restful.RouterParamas{
		BotnoiService: botnoiService,
	})

	//TODO: received message from kafka

	log.Println("start server on :", configs.Restful.Port)
	if err := fasthttp.ListenAndServe(":"+strconv.Itoa(configs.Restful.Port), restfulAPI.Handler); err != http.ErrServerClosed {
		log.Fatalf("main Error: failed on start server: %v", err)
	}
}
