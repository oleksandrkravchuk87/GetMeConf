FROM golang:1.9.2-alpine3.6 AS build

RUN mkdir -p /go/src \
&& mkdir -p /go/bin \
&& mkdir -p /go/pkg

ENV GOPATH=/go

ENV PATH=$GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/src/server \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/database \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/api


ADD ./server $GOPATH/src/server
ADD ./database $GOPATH/src/github.com/YAWAL/GetMeConf/database
ADD ./api $GOPATH/src/github.com/YAWAL/GetMeConf/api

ADD ./vendor $GOPATH/src/vendor
ADD ./Gopkg.lock $GOPATH/src/
ADD ./Gopkg.toml $GOPATH/src/

WORKDIR $GOPATH/src/server

RUN go build -o main .

CMD ["/go/src/server/main"]

EXPOSE 8081