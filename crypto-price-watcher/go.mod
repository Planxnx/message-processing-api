module github.com/Planxnx/message-processing-api/crypto-price-watcher

go 1.15

replace github.com/Planxnx/message-processing-api/message-schema => ../message-schema

require (
	github.com/Planxnx/message-processing-api/message-schema v0.0.0-20201214223757-1e1d612e1f29
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/ThreeDotsLabs/watermill-kafka/v2 v2.2.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/superoo7/go-gecko v1.0.0
	github.com/umbracle/go-web3 v0.0.0-20210510123804-6c88e3455435 // indirect
	google.golang.org/protobuf v1.26.0
)
