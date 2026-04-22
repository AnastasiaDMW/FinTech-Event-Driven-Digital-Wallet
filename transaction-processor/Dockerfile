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

RUN CGO_ENABLED=1 go build -o transaction-processor ./cmd/transaction-processor

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
    librdkafka1 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/transaction-processor .
COPY config/transactionprocessor.toml ./config/transactionprocessor.toml

EXPOSE 8081

CMD ["./transaction-processor"]