FROM golang:1.9.2-alpine3.6 AS build

RUN mkdir -p /go/src \
&& mkdir -p /go/bin \
&& mkdir -p /go/pkg

ENV GOPATH=/go

ENV PATH=$GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/src/restClient \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/database \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/api

ADD ./restClient $GOPATH/src/restClient
ADD ./database $GOPATH/src/github.com/YAWAL/GetMeConf/database
ADD ./api $GOPATH/src/github.com/YAWAL/GetMeConf/api

ADD ./vendor $GOPATH/src/vendor
ADD ./Gopkg.lock $GOPATH/src/
ADD ./Gopkg.toml $GOPATH/src/

WORKDIR $GOPATH/src/restClient

RUN go build -o main

CMD ["/go/src/restClient/main", "-port", "8080", "-serviceHost", "getmeconf_serverapp_1", "-servicePort", "3000"]

EXPOSE 8080