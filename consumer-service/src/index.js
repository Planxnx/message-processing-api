import kafka from "./connections/kafka";
const kafkaClient = kafka(process.env.KAFKA_HOST);
const consumerGroup1 = kafkaClient.consumer({ groupId: "group1" });

consumerGroup1.connect().then(async () => {
  await consumerGroup1.subscribe({
    topic: "message.new",
    fromBeginning: false,
  });
  await consumerGroup1.run({
    eachMessage: async ({ topic, partition, message }) => {
      const logMessage = JSON.stringify({
        status: `${process.env.SERVICE_NAME}: RECEIVED MESSAGE`,
        data: {
          topic,
          partition,
          offset: message.offset,
          message: message.value.toString(),
          time: message.timestamp,
        },
      });
      console.log(logMessage);
    },
  });
}).catch(err=>{
  console.log(`ERR: ConsumerGroup1Conenct: ${err}`)
});

console.log(`consumer service started - kafka`);
