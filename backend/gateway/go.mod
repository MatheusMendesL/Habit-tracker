module gateway

go 1.26.1

require (
	github.com/go-chi/chi/v5 v5.3.2
	go.uber.org/zap v1.18.1
	google.golang.org/grpc v1.79.3
	shared v0.0.0
)

replace shared => ../services/shared

require (
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
