FROM golang:1.9.2-alpine3.6 AS builder

RUN mkdir -p /go/src \
&& mkdir -p /go/bin \
&& mkdir -p /go/pkg

ENV GOPATH=/go

RUN mkdir -p $GOPATH/src/service \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/entitie \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/repository

ADD ./service $GOPATH/src/service
ADD entitie $GOPATH/src/github.com/YAWAL/GetMeConf/entitie
ADD ./repository $GOPATH/src/github.com/YAWAL/GetMeConf/repository

ADD ./vendor $GOPATH/src/vendor
ADD ./Gopkg.lock $GOPATH/src/
ADD ./Gopkg.toml $GOPATH/src/

WORKDIR $GOPATH/src/service

RUN go build -o $GOPATH/bin/service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app

WORKDIR /app

COPY --from=builder /go/bin/service .

CMD ["./service"]

EXPOSE $SERVICE_PORT