FROM golang:stretch AS build

RUN mkdir -p $GOPATH/src/github.com/astroband/astrologer
WORKDIR $GOPATH/src/github.com/astroband/astrologer

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w'

#===========================

FROM stellar/stellar-core:13.2.0-1297-b5dda51e AS stellar-core

#===========================

FROM stellar/base

ENV DATABASE_URL=postgres://localhost/core?sslmode=disable
ENV ES_URL=http://localhost:9200
ENV INGEST_GAP=-50

WORKDIR /root

COPY dependencies.sh entry.sh ./

RUN ["chmod", "+x", "./dependencies.sh"]
RUN ./dependencies.sh

COPY --from=stellar-core /usr/local/bin/stellar-core /usr/local/bin/

COPY --from=build /go/src/github.com/astroband/astrologer/astrologer .
RUN ["chmod", "+x", "./astrologer"]

ENTRYPOINT ["./entry.sh"]
CMD ./astrologer ingest -- $INGEST_GAP
