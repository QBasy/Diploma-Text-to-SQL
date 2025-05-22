FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./text-to-sql-service-go/cmd/main.go

EXPOSE 5006

CMD ["./server"]