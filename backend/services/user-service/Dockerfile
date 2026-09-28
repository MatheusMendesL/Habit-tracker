FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/services/user-service

RUN go mod tidy
RUN go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/services/user-service/app .
COPY --from=builder /app/services/user-service/cert ./cert

ENV SERVICE_NAME=user-service

EXPOSE 8080

CMD ["./app"]