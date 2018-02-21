FROM golang:1.9.2-alpine3.6 AS build

RUN mkdir -p /go/src \
&& mkdir -p /go/bin \
&& mkdir -p /go/pkg

ENV GOPATH=/go

ENV PATH=$GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/src/client \
&& mkdir -p $GOPATH/src/client/out \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/database \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/api

ADD ./client $GOPATH/src/client
ADD ./client/out $GOPATH/src/client/out
ADD ./database $GOPATH/src/github.com/YAWAL/GetMeConf/database
ADD ./api $GOPATH/src/github.com/YAWAL/GetMeConf/api

ADD ./vendor $GOPATH/src/vendor
ADD ./Gopkg.lock $GOPATH/src/
ADD ./Gopkg.toml $GOPATH/src/

WORKDIR $GOPATH/src/client

RUN go build -o main .

CMD ["/go/src/client/main", "-config-name", "mydom", "-config-type", "mongodb", "outpath", "$GOPATH/src/client/out"]

