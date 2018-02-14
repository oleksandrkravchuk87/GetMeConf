FROM golang:1.9.2-alpine3.6 AS build

RUN mkdir -p /go/src \
&& mkdir -p /go/bin \
&& mkdir -p /go/pkg

ENV GOPATH=/go

ENV PATH=$GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/src/client
ADD . $GOPATH/src/client

RUN mkdir -p $GOPATH/src/client/config
ADD ./config $GOPATH/src/client/config

RUN mkdir -p $GOPATH/src/client/config/in
ADD ./config/in $GOPATH/src/client/config/in

RUN mkdir -p $GOPATH/src/client/config/out

WORKDIR $GOPATH/src/client

RUN go build -o main .

CMD ["/go/src/client/main", "-config-id", "mongodb.json", "-config-path", "/go/src/client/config/in"]

