import * as line from '@line/bot-sdk';

import { Router } from 'express';
import config from '../config';
import mpaService from '../services/message-processing-api';

const router = Router();

//lazy cache
let userRefReplyToken = {};

const lineClient = new line.Client({
  channelAccessToken: config.line.channelAccessToken,
});

router.post('/message-processing-api', async (req, res, next) => {
  console.log(`Received MPA webhook events:${JSON.stringify(req.body)}\n`);
  let replyToken = userRefReplyToken[req.body.ref3];
  try {
    await lineClient.replyMessage(replyToken, {
      type: 'text',
      text: req.body.data.message,
    });
  } catch (error) {
    await lineClient.pushMessage(req.body.ref3, {
      type: 'text',
      text: req.body.data.message,
    });
  }
  res.json({});
});

router.post('/line', async (req, res, next) => {
  const events = req.body.events;
  console.log(`Received Line webhook events:${JSON.stringify(events)}\n`);
  for (let event of events) {
    if (event.type == 'message') {
      let userRef = event.source.userId;
      userRefReplyToken[userRef] = event.replyToken;
      mpaService.sendChitchat(event.message.text, userRef);
    }
  }
  res.json({});
});

export default router;
