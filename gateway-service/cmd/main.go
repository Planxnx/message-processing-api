package main

import (
	"log"

	"github.com/Planxnx/message-processing-api/gateway-service/api"
	"github.com/valyala/fasthttp"
)

var (
	port string = "8080"
)

func main() {

	routerDependency := &api.RouterDependency{}
	router := routerDependency.GetFastHTTPRouter().Handler

	log.Println("start server on :" + port)
	log.Fatal(fasthttp.ListenAndServe(":"+port, router))
}
