FROM golang:alpine AS build

RUN apk add --no-cache git

RUN mkdir -p $GOPATH/src/github.com/astroband/astrologer
WORKDIR $GOPATH/src/github.com/astroband/astrologer

ADD . .

RUN GO111MODULE=on go build

# ===============================================================================================

FROM alpine:latest

ENV DATABASE_URL=postgres://localhost/core?sslmode=disable
ENV ES_URL=http://localhost:9200
ENV INGEST_GAP=-50

WORKDIR /root

COPY --from=build /go/src/github.com/astroband/astrologer/astrologer .
RUN chmod +x ./astrologer

CMD /root/astrologer ingest -- $INGEST_GAP