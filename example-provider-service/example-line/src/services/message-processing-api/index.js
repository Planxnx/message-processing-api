import axios from 'axios';
import config from '../../config';

const {
  url: messageProcessingAPI,
  id: providerID,
  token: authToken,
} = config.mpa;

export default {
  sendChitchat: async function (message, userRef) {
    const options = {
      method: 'POST',
      url: messageProcessingAPI,
      headers: {
        Authorization: authToken,
        'Provider-ID': providerID,
        'Content-Type': 'application/json',
      },
      data: {
        message: message,
        userRef: userRef,
        features: {
          chitchat: true,
        },
      },
    };
    try {
      const resp = await axios(options);
      return resp.data;
    } catch (error) {}
  },
};
