FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/habit-service

RUN go mod tidy
RUN go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/habit-service/app .

EXPOSE 8082

CMD ["./app"]