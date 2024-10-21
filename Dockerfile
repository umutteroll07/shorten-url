FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /shorten-url/cmd/main ./cmd

EXPOSE 3000

CMD ["/shorten-url/cmd/main"]