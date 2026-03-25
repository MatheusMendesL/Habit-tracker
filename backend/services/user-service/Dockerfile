FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

# entra no serviço
WORKDIR /app/user-service

RUN go mod tidy
RUN go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/user-service/app .

EXPOSE 8080

CMD ["./app"]