package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Planxnx/message-processing-api/crypto-price-watcher/internal/coingecko"
	"github.com/Planxnx/message-processing-api/crypto-price-watcher/internal/messageapi"
)

const (
	ServiceName           = "crypto-price-watcher"
	FeatureCryptoGetPrice = "crypto-get-price"
)

func main() {
	ctx := context.Background()

	//Initial Dependency
	messageClient, err := messageapi.New()
	if err != nil {
		log.Fatalf("main Error: failed on create message servic: %v\n", err)
	}
	coingeckoClient, err := coingecko.New()
	if err != nil {
		log.Fatalf("main Error: failed on create coingecko client: %v\n", err)
	}

	//Define Feature Config
	cryptoGetPriceConfig := &messageapi.HandlerConfig{
		FeatureName:      FeatureCryptoGetPrice,
		ServiceName:      ServiceName,
		AsynchronousMode: true,
		SynchronousMode:  true,
	}

	//Add Feature Handler
	messageClient.AddHandler(messageschema.CommonMessageTopic, cryptoGetPriceConfig, CryptoGetPriceHandler(coingeckoClient))

	//Start service
	log.Printf("%s: Start messagequeue subscriber ...\n", ServiceName)
	if err := messageClient.Run(ctx); err != nil {
		log.Fatalf("main Error: failed on start messagequeue subscruber: %v", err)
	}
}

//CryptoGetPriceHandler message api handler for crypto get price
func CryptoGetPriceHandler(coingeckoClient *coingecko.CoinGecko) messageapi.MessageAPIHandler {

	type CryptoGetPriceReqData struct {
		CoinSymbol string `json:"coin,omitempty"`
	}

	return func(replymessage, request *messageschema.DefaultMessage) error {

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
}
