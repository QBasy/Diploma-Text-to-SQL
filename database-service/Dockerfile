FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o database-service

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/database-service .

EXPOSE 5002

CMD ["./database-service"]