FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

RUN go build -o server ./cmd/server

EXPOSE 5006

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./server"]
