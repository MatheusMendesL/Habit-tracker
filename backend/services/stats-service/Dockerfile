FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/services/stats-service

RUN go mod tidy
RUN go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/services/stats-service/app .
COPY --from=builder /app/services/stats-service/cert ./cert

ENV SERVICE_NAME=stats-service

EXPOSE 8083

CMD ["./app"]