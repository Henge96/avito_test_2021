FROM golang:1.18

WORKDIR /app

COPY . /app

RUN go build ./cmd/main.go

EXPOSE 8080

CMD ./main -config ./cmd/config.yml

