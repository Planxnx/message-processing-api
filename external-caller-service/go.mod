module github.com/Planxnx/message-processing-api/external-caller-service

go 1.14

replace github.com/Planxnx/message-processing-api/message-schema => ../message-schema

require (
	github.com/Planxnx/message-processing-api/message-schema v0.0.0-20201119085137-32a2c81d2015
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/ThreeDotsLabs/watermill-kafka/v2 v2.2.0
	github.com/buaazp/fasthttprouter v0.1.1
	github.com/google/uuid v1.1.2 // indirect
	github.com/klauspost/compress v1.11.0 // indirect
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/valyala/fasthttp v1.16.0
	golang.org/x/sys v0.0.0-20200909081042-eff7692f9009 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/protobuf v1.25.0
)
