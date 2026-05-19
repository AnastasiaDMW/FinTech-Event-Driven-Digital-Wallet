#!/bin/bash
set -e

echo "STARTING KAFKA CONNECT..."

/etc/confluent/docker/run &

CONNECT_PID=$!

CONNECT_URL="http://localhost:8083"

echo "WAITING FOR KAFKA CONNECT REST API..."

until curl -s "$CONNECT_URL/" | grep -q "version"; do
  echo "Kafka Connect not ready yet..."
  sleep 5
done

echo "KAFKA CONNECT READY"

for file in /connectors/*.json; do
  name=$(basename "$file" .json)

  echo "REGISTERING CONNECTOR: $name"

  tmp="/tmp/$name.json"

  envsubst '
  ${DEBEZIUM_USER}
  ${DEBEZIUM_PASSWORD}
  ${BALANCE_NAME}
  ${BALANCE_DBNAME}
  ${BALANCE_HOSTNAME}
  ${BALANCE_SERVER_NAME}
  ${TRANSACTIONAL_ORCHESTRATOR_NAME}
  ${TRANSACTIONAL_ORCHESTRATOR_DBNAME}
  ${TRANSACTIONAL_ORCHESTRATOR_HOSTNAME}
  ${TRANSACTIONAL_ORCHESTRATOR_SERVER_NAME}
  ${MONGO_CONNECTION_URI}
  ' < "$file" > "$tmp"

  echo "CONNECTOR CONFIG:"
  cat "$tmp"

  response=$(curl -s \
    -o /tmp/resp.json \
    -w "%{http_code}" \
    -X POST "$CONNECT_URL/connectors" \
    -H "Content-Type: application/json" \
    -d @"$tmp")

  echo "HTTP STATUS: $response"

  cat /tmp/resp.json

  if [ "$response" != "201" ] && [ "$response" != "409" ]; then
    echo "FAILED TO REGISTER CONNECTOR: $name"
    exit 1
  fi
done

echo "ALL CONNECTORS REGISTERED"

wait $CONNECT_PID