const db = db.getSiblingDB("logs_db");

db.createUser({
  user: "kafka_writer",
  pwd: "secure_password",
  roles: [
    { role: "readWrite", db: "logs_db" }
  ]
});

db.createCollection("application_logs");

db.application_logs.createIndex(
  { timestamp: 1 },
  { expireAfterSeconds: 604800 }
);

db.application_logs.createIndex({ service: 1, timestamp: -1 });
db.application_logs.createIndex({ level: 1 });