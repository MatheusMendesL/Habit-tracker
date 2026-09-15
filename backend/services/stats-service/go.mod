module stats-service

go 1.26

replace shared => ../shared

require (
	github.com/jackc/pgx/v5 v5.10.0
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.79.3
	shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

require (
	github.com/google/uuid v1.6.0
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	go.uber.org/zap v1.28.0
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)
