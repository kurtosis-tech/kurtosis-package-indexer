module github.com/kurtosis-tech/kurtosis-package-indexer/server

go 1.19

replace github.com/kurtosis-tech/kurtosis-package-indexer/api/golang => ../api/golang

require (
	github.com/kurtosis-tech/kurtosis-package-indexer/api/golang v0.0.0 // local dependency
	github.com/kurtosis-tech/minimal-grpc-server/golang v0.0.0-20230710164206-90b674acb269
	github.com/kurtosis-tech/stacktrace v0.0.0-20211028211901-1c67a77b5409
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.57.0
)

require (
	github.com/rs/cors v1.9.0
	google.golang.org/protobuf v1.31.0
)

require (
	connectrpc.com/connect v1.11.0 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/improbable-eng/grpc-web v0.15.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/kurtosis-tech/kurtosis/connect-server v0.0.0-20230825003324-75d481e0db8c // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230803162519-f966b187b2e5 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)
