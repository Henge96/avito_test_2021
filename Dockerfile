FROM golang:1.18.1-alpine3.15

WORKDIR /app

COPY . /app

RUN go build ./cmd/main.go

EXPOSE 8080

CMD ./main -config ./cmd/config.yml

