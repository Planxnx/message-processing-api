import * as line from '@line/bot-sdk';

import {
  Router
} from 'express';
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
  console.log(`Got Reply Token ${replyToken}`)
  try {
    await lineClient.replyMessage(replyToken, {
      type: 'text',
      text: req.body.data.message ? req.body.data.message : JSON.stringify(req.body.data),
    });
  } catch (error) {
    console.log(JSON.stringify(error))
    await lineClient.pushMessage(req.body.ref3, {
      type: 'text',
      text: "เกิดข้อผิดพลาดบางอย่างขึ้น",
    });
  }
  res.json({});
});

router.post('/line', async (req, res, next) => {
  const events = req.body.events;
  console.log(`Received Line webhook events:${JSON.stringify(events)}\n`);
  try {
    for (let event of events) {
      if (event.type == 'message') {
        const message = event.message.text;
        if (message.startsWith('/')) {
          if (message == '/help') {
            lineClient.replyMessage(event.replyToken, {
              type: 'text',
              text: 'HELP🚑\nช่วยตัวเองก่อนน้า\n',
            });
          } else if (message.startsWith('/lotto')) {
            let userRef = event.source.userId;
            let lottoID = message.slice(7);
            mpaService.checkLottoReward(lottoID, userRef).then((data) => {
              let result = '';
              if (data.foundReward) {
                result = `ยินดีด้วย หมายเลข ${data.foundReward[0].number} ได้รางวัล ${data.foundReward[0].name} เป็นเงินมูลค่า ${data.foundReward[0].reward} บาท`;
              } else {
                result = `เสียใจด้วย หมายเลข ${lottoID} ไม่ถูกรางวัล`;
              }
              lineClient.replyMessage(event.replyToken, {
                type: 'text',
                text: result,
              });
            });
          } else if (message.startsWith('/crypto-price-async')) {
            let userRef = event.source.userId;
            userRefReplyToken[userRef] = event.replyToken;
            let coinSymbol = String(message.slice(20)).toUpperCase();
            mpaService.checkCryptoPrice(userRef, coinSymbol, true)
          } else if (message.startsWith('/crypto-price')) {
            let userRef = event.source.userId;
            let coinSymbol = String(message.slice(14)).toUpperCase();
            mpaService.checkCryptoPrice(userRef, coinSymbol, false).then(({
              data: coins
            }) => {
              let result = '';
              for (let name in coins) {
                result += `💰เหรียญ ${name}\n\t1 ${coinSymbol} = ${coins[name].usd} USD\n`
              }
              lineClient.replyMessage(event.replyToken, {
                type: 'text',
                text: result.length > 0 ? result : `ขออภัย ไม่เจอเหรียญ ${coinSymbol} ในระบบ`,
              });
            });
          } else {
            lineClient.replyMessage(event.replyToken, {
              type: 'text',
              text: "ดูคำสั่งเพิ่มเติ่ม พิมพ์ '/help'",
            });
          }
          continue;
        }
        let userRef = event.source.userId;
        userRefReplyToken[userRef] = event.replyToken;
        mpaService.sendChitchat(message, userRef);
      }
    }
  } catch (error) {
    console.log(`Error Line webhook events:${JSON.stringify(events)}\n`);
  }
  res.json({});
});

export default router;