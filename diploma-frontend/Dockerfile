FROM node:20

WORKDIR /app

COPY package*.json ./

# Install dependencies
RUN npm install

COPY . .

RUN npm run dev

EXPOSE 5173

CMD ["npm", "run", "preview"]
