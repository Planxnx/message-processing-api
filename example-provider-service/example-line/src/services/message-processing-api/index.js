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
        feature: 'Chitchat',
      },
    };
    try {
      const resp = await axios(options);
      return resp.data;
    } catch (error) {
      console.log(error);
    }
  },
  checkLottoReward: async function (lottoID, userRef) {
    const options = {
      method: 'POST',
      url: messageProcessingAPI + '/sync',
      headers: {
        Authorization: authToken,
        'Provider-ID': providerID,
        'Content-Type': 'application/json',
      },
      data: {
        userRef: userRef,
        feature: 'Check-Latest-Lottery',
        data: {
          number: lottoID,
        },
      },
    };
    try {
      const resp = await axios(options);
      return resp.data;
    } catch (error) {
      console.log(error);
    }
  },
  checkCryptoPrice: async function (userRef, coinSymmol, isAsync) {
    const options = {
      method: 'POST',
      url: `${messageProcessingAPI}${isAsync ? "" : "/sync"}`,
      headers: {
        Authorization: authToken,
        'Provider-ID': providerID,
        'Content-Type': 'application/json',
      },
      data: {
        userRef: userRef,
        feature: 'crypto-get-price',
        data: {
          coin: coinSymmol,
        },
      },
    };
    try {
      const resp = await axios(options);
      return resp.data;
    } catch (error) {
      console.log(error);
    }
  },
};