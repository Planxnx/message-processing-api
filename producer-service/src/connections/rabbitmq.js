import amqp from "amqplib";

const createRabbitMQConnection = (rabbitHost) => {
  const rabbitMQHost =
    rabbitHost || process.env.RABBITMQ_HOST || "localhost:5672";
  return amqp.connect(`amqp://admin:admin@${rabbitMQHost}`);
};

const createRabbitMQChannel = async (rabbitHost) => {
  const rabbitMQHost =
    rabbitHost || process.env.RABBITMQ_HOST || "localhost:5672";
  const rabbitMQConnection = await amqp.connect(
    `amqp://admin:admin@${rabbitMQHost}`
  );
  return rabbitMQConnection.createChannel();
};

export { createRabbitMQConnection, createRabbitMQChannel };
