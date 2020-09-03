import { Kafka } from "kafkajs";

export default (host) => {
  const kafkaHost = host || process.env.KAFKA_HOST || "localhost:9092";
  const kafkaClient = new Kafka({
    clientId: "producer-service",
    brokers: [kafkaHost],
  });
  console.log(`kafka host at - ${kafkaHost}`);
  return kafkaClient;
};
