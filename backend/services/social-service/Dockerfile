FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/services/social-service

RUN go mod tidy
RUN go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/services/social-service/app .
COPY --from=builder /app/services/social-service/cert ./cert

ENV SERVICE_NAME=social-service

EXPOSE 8081

CMD ["./app"]