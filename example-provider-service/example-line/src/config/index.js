import * as dotenv from 'dotenv';

dotenv.config();

export default {
  mpa: {
    url: process.env.MPA_API,
    id: process.env.MPA_PROVIDER_ID,
    token: process.env.MPA_TOKEN,
  },
  line: {
    channelId: process.env.LINE_CHANNEL_ID,
    channelSecret: process.env.LINE_CHANNEL_SECRET,
    channelAccessToken: process.env.LINE_ACCESS_TOKEN,
  },
};