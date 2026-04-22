#!/bin/bash
set -e

CONNECT_URL="http://kafka-connect:8083"

until curl -s "$CONNECT_URL/" | grep -q "version"; do
  echo "Kafka Connect ПРОВЕРКА..."
  sleep 1
done

echo "Kafka Connect ГОТОВА"

for file in /connectors/*.json; do
  name=$(basename "$file" .json)
  echo "Загрузка $name"
  tmp="/tmp/$name.json"
envsubst '${DEBEZIUM_USER} ${DEBEZIUM_PASSWORD} ${BALANCE_NAME} ${BALANCE_DBNAME} ${BALANCE_HOSTNAME} ${BALANCE_SERVER_NAME} ${TRANSACTIONAL_ORCHESTRATOR_NAME} ${TRANSACTIONAL_ORCHESTRATOR_DBNAME} ${TRANSACTIONAL_ORCHESTRATOR_HOSTNAME} ${TRANSACTIONAL_ORCHESTRATOR_SERVER_NAME} ${MONGO_CONNECTION_URI}' < "$file" > "$tmp"
  cat "$tmp"
  response=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -X POST "$CONNECT_URL/connectors" \
    -H "Content-Type: application/json" \
    -d @"$tmp")
  echo "HTTP status: $response"
  cat /tmp/resp.json

  if [ "$response" != "201" ] && [ "$response" != "409" ]; then
    echo "ОШИБКА $name"
    exit 1
  fi
done
echo "ВСЕ КОННЕКТОРЫ ВЫПОЛНЕНЫ"