module github.com/kurtosis-tech/kurtosis-package-indexer/server

go 1.19

replace github.com/kurtosis-tech/kurtosis-package-indexer/api/golang => ../api/golang

require (
	connectrpc.com/connect v1.11.0
	github.com/aws/aws-sdk-go v1.44.334
	github.com/google/go-github/v54 v54.0.0
	github.com/kurtosis-tech/kurtosis-package-indexer/api/golang v0.0.0-00010101000000-000000000000
	github.com/kurtosis-tech/kurtosis/connect-server v0.0.0-20230825003324-75d481e0db8c
	github.com/kurtosis-tech/stacktrace v0.0.0-20211028211901-1c67a77b5409
	github.com/rs/cors v1.9.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
	github.com/tilt-dev/starlark-lsp v0.0.0-20230210155543-84c13fe0ff94
	go.etcd.io/bbolt v1.3.7
	go.etcd.io/etcd/client/v3 v3.5.9
	go.starlark.net v0.0.0-20230814145427-12f4cb8177e4
	golang.org/x/oauth2 v0.11.0
	google.golang.org/protobuf v1.31.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/ProtonMail/go-crypto v0.0.0-20230217124315-7d5c6f04bbb8 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.etcd.io/etcd/api/v3 v3.5.9 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.9 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.20.0 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230726155614-23370e0ffb3e // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230706204954-ccb25ca9f130 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230803162519-f966b187b2e5 // indirect
	google.golang.org/grpc v1.57.0 // indirect
)
