FROM golang:1.26 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/backend/services/habit-service

RUN go mod tidy
RUN go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/backend/services/habit-service/app .
COPY --from=builder /app/backend/services/habit-service/cert ./cert

ENV SERVICE_NAME=habit-service

EXPOSE 8082/tcp

CMD ["./app"]