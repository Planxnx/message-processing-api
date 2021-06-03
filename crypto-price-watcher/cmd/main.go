package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Planxnx/message-processing-api/crypto-price-watcher/internal/message"
	"github.com/Planxnx/message-processing-api/crypto-price-watcher/pkg/coingecko"
	kafapkg "github.com/Planxnx/message-processing-api/crypto-price-watcher/pkg/kafka"
	messageSchema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/robfig/cron/v3"
)

const (
	ServiceName           = "crypto-price-watcher"
	FeatureCryptoGetPrice = "crypto-get-price"
)

func main() {
	ctx := context.Background()

	//Init Handler Config
	CryptoGetPriceHandlerConfig := &message.Config{
		FeatureName:      FeatureCryptoGetPrice,
		ServiceName:      ServiceName,
		AsynchronousMode: true,
		SynchronousMode:  true,
	}

	//Initial Dependency
	kafkaSubscriber, err := kafapkg.NewSubscriber()
	if err != nil {
		log.Fatalf("main Error: failed on create kafka subscriber: %v", err)
	}
	kafkaNewPublisher, err := kafapkg.NewPubliser()
	if err != nil {
		log.Fatalf("main Error: failed on create kafka publisher: %v", err)
	}
	coingeckoClient, err := coingecko.New()
	if err != nil {
		log.Fatalf("main Error: failed on create coingecko client: %v", err)
	}
	messageUsecase := message.NewUsecase(kafkaNewPublisher)

	//Create Handler
	CryptoGetPriceHandler := func(replymessage, request *messageSchema.DefaultMessage) error {
		type CryptoGetPriceReqData struct {
			CoinSymbol string `json:"coin,omitempty"`
		}
		requestData := &CryptoGetPriceReqData{}

		if err := json.Unmarshal(request.Data, requestData); err != nil {
			replymessage.ErrorInternal = err.Error()
			replymessage.PublishedAt = timestamppb.Now()
			log.Printf("Handler ERROR:  Error: failed on unmarshal message: %v\n", err)
			return err
		}

		if len(requestData.CoinSymbol) == 0 {
			replymessage.Error = "data.coin is requried"
			replymessage.PublishedAt = timestamppb.Now()
			return errors.New(replymessage.Error)
		}

		pricesData, err := coingeckoClient.GetCoinPrices(requestData.CoinSymbol)
		if err != nil {
			replymessage.ErrorInternal = err.Error()
			replymessage.PublishedAt = timestamppb.Now()
			log.Printf("Handler ERROR:  Error: failed on get coin prices: %v\n", err)
			return err
		}

		attchData, err := json.Marshal(pricesData)
		if err != nil {
			replymessage.ErrorInternal = err.Error()
			replymessage.PublishedAt = timestamppb.Now()
			log.Printf("Handler ERROR:  Error: failed on json marshal: %v\n", err)
			return err
		}

		replymessage.Data = attchData
		replymessage.PublishedAt = timestamppb.Now()
		replymessage.PublishedAt = timestamppb.Now()

		return nil
	}

	//Publish Status to APIGatway
	go healthCheck(ServiceName, kafkaNewPublisher)

	//Create Subscriber Client
	messages, err := kafkaSubscriber.Subscribe(ctx, messageSchema.CommonMessageTopic)
	if err != nil {
		log.Fatalf("main Error: failed on subscribe topic: %v", err)
	}

	log.Printf("%s: Start messagequeue subscriber ...\n", ServiceName)
	//Start CryptoGetPrice Subscriber Message
	for msg := range messages {
		messageUsecase.CreateMessageHandler(CryptoGetPriceHandlerConfig, CryptoGetPriceHandler)(msg)
	}
}

func healthCheck(serviceName string, kafkaPublisher *kafka.Publisher) {
	healthCheckCmd := func() {

		//Chitchat HealthCheck
		go func() {
			chitchat := &messageschema.HealthCheckMessage{
				Feature:     "crypto-get-price",
				Description: "ตรวจสอบราคาล่าสุดของคริปโต",
				ExecuteMode: []messageschema.ExecuteMode{
					messageschema.ExecuteMode_Asynchronous,
					messageschema.ExecuteMode_Synchronous,
				},
				ServiceName: serviceName,
			}

			chitchatByte, err := proto.Marshal(chitchat)
			if err != nil {
				log.Println("health check error: can't marshal crypto-get-price message")
			}

			//Publish status
			if err := kafkaPublisher.Publish(messageschema.HealthCheckTopic, watermillmessage.NewMessage(watermill.NewShortUUID(), chitchatByte)); err != nil {
				log.Printf("health check error: failed on publish crypto-get-price message: %v\n", err)
			}
		}()

	}

	//Startting
	go healthCheckCmd()

	c := cron.New()
	c.AddFunc("@every 3m", healthCheckCmd)
	c.Start()
}
