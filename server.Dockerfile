FROM golang:1.9.2-alpine3.6 AS build

RUN mkdir -p /go/src \
&& mkdir -p /go/bin \
&& mkdir -p /go/pkg

ENV GOPATH=/go

ENV PATH=$GOPATH/bin:$PATH

ENV PORT=3000

ENV PDB_HOST=getmeconf_db_1
ENV PDB_PORT=5432
ENV PDB_USER=postgres
ENV PDB_PASSWORD=root
ENV PDB_NAME=postgres
ENV MAX_OPENED_CONNECTIONS_TO_DB=5
ENV MAX_IDLE_CONNECTIONS_TO_DB=0
ENV MB_CONN_MAX_LIFETIME_MINUTES=30

ENV CACHE_EXPIRATION_TIME=5
ENV CACHE_CLEANUP_INTERVAL=10

RUN mkdir -p $GOPATH/src/service \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/repository \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/api


ADD ./server $GOPATH/src/service
ADD ./database $GOPATH/src/github.com/YAWAL/GetMeConf/repository
ADD ./api $GOPATH/src/github.com/YAWAL/GetMeConf/api

ADD ./vendor $GOPATH/src/vendor
ADD ./Gopkg.lock $GOPATH/src/
ADD ./Gopkg.toml $GOPATH/src/

WORKDIR $GOPATH/src/server

RUN go build -o main .

CMD ["/go/src/server/main"]

EXPOSE $PORT