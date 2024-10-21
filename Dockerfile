FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go.mod download

COPY ..

RUN go build -o /url-shortener/cmd/server/main.go ./cmd/server

EXPOSE 3000

CMD ["./main"]