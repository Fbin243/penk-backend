FROM node:22-alpine AS builder

WORKDIR /tenk-backend

COPY services/gateway/package*.json .

RUN npm install

COPY services/gateway .

CMD ["sh", "-c", "npm run build && npm run start"]
