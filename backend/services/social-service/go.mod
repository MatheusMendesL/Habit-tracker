module social

go 1.25.4

replace shared => ../shared

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/joho/godotenv v1.5.1
	go.uber.org/zap v1.27.1
	google.golang.org/grpc v1.79.3
	shared v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
