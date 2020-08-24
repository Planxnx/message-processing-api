import express from "express";
import bodyParser from "body-parser";
import kafka from "./connections/kafka";
const kafkaClient = kafka();

const app = express();

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

app.get("/", (req, res) => {
  res.status(200).send({ status: 200, message: "ok" });
});

app.post("/message", async (req, res) => {
  try {
    const { message, topic } = req.body;
    const producer = kafkaClient.producer();
    await producer.connect();
    await producer.send({
      topic: topic || "TEST_TOPIC",
      messages: [
        {
          value: message || "TEST_MESSAGE",
        },
      ],
    });
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
    });
  }
});

const port = process.env.SERVER_PORT || 8080;
app.listen(port, () => {
  console.log(`producer service started - ${port}`);
});
