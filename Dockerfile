FROM golang:1.9.2-alpine3.6 AS build

RUN mkdir -p /go/src \
&& mkdir -p /go/bin \
&& mkdir -p /go/pkg

ENV GOPATH=/go

ENV PATH=$GOPATH/bin:$PATH

#ENV PDB_HOST=getmeconf_db_1
#ENV PDB_PORT=5432
#ENV PDB_USER=postgres
#ENV PDB_PASSWORD=root
#ENV PDB_NAME=postgres

RUN mkdir -p $GOPATH/src/service \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/entities \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/repository \
&& mkdir -p $GOPATH/src/github.com/YAWAL/GetMeConf/api


ADD ./service $GOPATH/src/service
ADD entities $GOPATH/src/github.com/YAWAL/GetMeConf/entities
ADD ./repository $GOPATH/src/github.com/YAWAL/GetMeConf/repository
ADD ./api $GOPATH/src/github.com/YAWAL/GetMeConf/api

ADD ./vendor $GOPATH/src/vendor
ADD ./Gopkg.lock $GOPATH/src/
ADD ./Gopkg.toml $GOPATH/src/

WORKDIR $GOPATH/src/service

RUN go build -o $GOPATH/bin/service .

RUN rm -rf /GOPATH/src && rm -rf /GOPATH/pkg

CMD ["/go/bin/service"]

EXPOSE $PORT