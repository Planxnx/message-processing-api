import express from "express";
import bodyParser from "body-parser";
import { createRabbitMQChannel } from "./connections/rabbitmq";

const app = express();

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

app.get("/", (req, res) => {
  res.status(200).send({ status: 200, message: "ok" });
});

app.post("/message", async (req, res) => {
  try {
    const { message, topic } = req.body;
    const mqchannel = await createRabbitMQChannel();
    await mqchannel.assertQueue(topic || "TEST_TOPIC");
    mqchannel.sendToQueue(
      topic || "TEST_TOPIC",
      Buffer.from(message || "TEST_MESSAGE")
    );
    res.status(200).send({
      status: 200,
      message: "produce success",
      data: {
        topic,
        message,
      },
    });
  } catch (error) {
    res.status(400).send({
      status: 400,
      message: `produce failed: ${error}`,
      data: {
        topic,
        message,
      },
    });
  }
});

const port = process.env.MOCK_SERVER_PORT || 8888;
app.listen(port, () => {
  console.log(`producer service started - ${port}`);
});
