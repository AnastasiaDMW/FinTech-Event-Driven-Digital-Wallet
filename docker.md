Комманды для полного запуск:
docker-compose -f docker-compose-kafka.yml -f docker-compose.yml up -d

Команды для запуска кафки с интерфейсом
docker-compose -f docker-compose-database.yml -f docker-compose-kafka.yml up -d

Bootstrap.server "localhost:19092,localhost:29092,localhost:39092"

загрузка в кафка коннектор
curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json"  -d @Debezium.json