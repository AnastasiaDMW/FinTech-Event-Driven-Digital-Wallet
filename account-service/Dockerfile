FROM golang:1.26 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y \
    git \
    gcc \
    libc6-dev \
    librdkafka-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o account-service ./cmd/account-service

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
    librdkafka1 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/account-service .
COPY ./.env ./.env
COPY config/accountservice.toml ./config/accountservice.toml

EXPOSE 8081

CMD ["./account-service"]