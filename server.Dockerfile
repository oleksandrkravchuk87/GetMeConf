FROM golang:1.9.2-alpine3.6 AS build

RUN mkdir -p /go/src \
&& mkdir -p /go/bin \
&& mkdir -p /go/pkg

ENV GOPATH=/go

ENV PATH=$GOPATH/bin:$PATH

ENV PORT=3000

ENV DB_HOST=getmeconf_db_1
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=root
ENV DB_NAME=postgres
ENV MAX_OPENED_CONNECTIONS_TO_DB=5
ENV MAX_IDLE_CONNECTIONS_TO_DB=0
ENV MB_CONN_MAX_LIFETIME_MINUTES=30

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

EXPOSE $PORT