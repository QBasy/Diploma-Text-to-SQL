FROM node:20 AS builder

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

FROM node:20 AS production

WORKDIR /app

COPY --from=builder /app /app

EXPOSE 4173

CMD ["npm", "run", "preview", "--", "--host", "0.0.0.0"]