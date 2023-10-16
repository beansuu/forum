FROM golang:latest

WORKDIR /app

COPY . .

RUN go get github.com/mattn/go-sqlite3

RUN go build -o main ./cmd/server

CMD ["./main"]