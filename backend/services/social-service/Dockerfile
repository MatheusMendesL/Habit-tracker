FROM golang:1.26 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/backend/services/social-service

RUN go mod tidy
RUN go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/backend/services/social-service/app .
COPY --from=builder /app/backend/services/social-service/cert ./cert

ENV SERVICE_NAME=social-service

EXPOSE 8081/tcp

CMD ["./app"]