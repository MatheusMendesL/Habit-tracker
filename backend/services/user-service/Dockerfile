FROM golang:1.26 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/backend/services/user-service

RUN go mod tidy
RUN go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/backend/services/user-service/app .
COPY --from=builder /app/backend/services/user-service/cert ./cert

ENV SERVICE_NAME=user-service

EXPOSE 8080/tcp

CMD ["./app"]