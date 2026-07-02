FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/social-service

RUN go mod tidy
RUN go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/social-service/app .

EXPOSE 8081

CMD ["./app"]