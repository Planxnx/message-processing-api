FROM node:14.8.0-stretch-slim

ENV SERVICE_PORT=4000

WORKDIR /usr/app

COPY ./ ./ 

RUN npm install

EXPOSE $SERVICE_PORT

CMD ["npm","start"]