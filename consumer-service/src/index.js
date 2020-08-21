import { createRabbitMQChannel } from "./connections/rabbitmq";

createRabbitMQChannel().then((mqchannel) => {
  mqchannel.consume(
    "NEW_MESSAGE",
    (msg) => {
      const logsMessage = JSON.stringify({
        status: "RECEIVED MESSAGE",
        data: {
          topic: msg.fields.routingKey,
          message: msg.content,
        },
      });
      console.log(logsMessage);
    },
    { noAck: false }
  );
});

console.log(`consumer service started - rabbitmq`);
