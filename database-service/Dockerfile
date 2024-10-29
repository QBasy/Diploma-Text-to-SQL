# Start from the official Go image for development
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o database-service .

EXPOSE 5002

CMD ["./main"]
