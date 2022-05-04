FROM golang:1.13.8

WORKDIR /app

COPY . /app

RUN go build ./cmd/main.go

EXPOSE 8080

CMD ./main