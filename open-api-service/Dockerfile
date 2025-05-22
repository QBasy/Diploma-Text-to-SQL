FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./open-api-service/cmd/server

EXPOSE 5009

CMD ["./server"]