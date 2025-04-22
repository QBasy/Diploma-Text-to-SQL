FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./text-to-sql/cmd/

EXPOSE 5006

CMD ["./server"]