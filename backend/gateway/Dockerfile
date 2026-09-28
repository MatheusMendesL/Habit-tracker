FROM golang:1.26 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/gateway

RUN go mod tidy
RUN go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/gateway/app .
COPY --from=builder /app/services/user-service/cert/ca-cert.pem ./cert/ca-cert.pem

ENV SERVICE_NAME=gateway
ENV TLS_CA_CERT_PATH=/app/cert/ca-cert.pem
ENV PORT=:8084

EXPOSE 8084/tcp

CMD ["./app"]
