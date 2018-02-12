FROM golang:1.9.2-alpine3.6 AS build

WORKDIR /server

ADD . /server

RUN go build -o main .

CMD ["/server/main"]

EXPOSE 8081