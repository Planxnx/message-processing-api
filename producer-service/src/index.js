import express from "express";
import bodyParser from "body-parser";
import kafka from "./connections/kafka";
const kafkaClient = kafka(process.env.KAFKA_HOST);

const app = express();

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

app.post("/", (req, res) => {
  console.log(`${process.env.SERVICE_NAME}: health ok`);
  console.log(`received:  ${JSON.stringify(req.body)}`)
  res.status(200).send({ status: 200, message: "ok" });
});

app.get("/health", (req, res) => {
  console.log(`${process.env.SERVICE_NAME}: health ok`);
  res.status(200).send({ status: 200, message: "ok" });
});

app.get("/kafka-health", async (req, res) => {
  try {
    const kafkaAdmin = kafkaClient.admin();
    await kafkaAdmin.connect();
    console.log(`${process.env.SERVICE_NAME}: kafka health ok`);
    res.status(200).send({ status: 200, message: "kafka ok" });
  } catch (error) {
    console.log(`${process.env.SERVICE_NAME} ERR: kafka error ${error}`);
    res.status(500).send({ status: 500, message: error });
  }
});
app.post("/kafka-health", async (req, res) => {
  try {
    const producer = kafkaClient.producer();
    await producer.connect();;
    await producer.send({
      topic: "message.test",
      messages: [
        {
          value: req.body.message || "TEST_MESSAGE",
        },
      ],
    });
    console.log(`${process.env.SERVICE_NAME}: kafka health ok`);
    res.status(200).send({ status: 200, message: "kafka ok" });
  } catch (error) {
    console.log(`${process.env.SERVICE_NAME} ERR: kafka error ${error}`);
    res.status(500).send({ status: 500, message: error });
  }
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
    console.log(`${process.env.SERVICE_NAME}: SENT`);
    res.status(200).send({
      status: 200,
      message: "produce success",
      data: {
        topic,
        message,
      },
    });
  } catch (error) {
    console.log(`${process.env.SERVICE_NAME} ERR: produer error ${error}`);
    res.status(400).send({
      status: 400,
      message: `produce failed: ${error}`,
    });
  }
});

const port = process.env.SERVICE_PORT || 8080;
app.listen(port, () => {
  console.log(`${process.env.SERVICE_NAME} service started - ${port}`);
});
