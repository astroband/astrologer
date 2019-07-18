FROM golang:alpine AS build

RUN apk add --no-cache git curl

RUN mkdir -p $GOPATH/src/github.com/astroband
RUN cd $GOPATH/src/github.com/astroband && git clone https://github.com/astroband/astrologer.git

WORKDIR $GOPATH/src/github.com/astroband/astrologer

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure -v
RUN go build -v

# ===============================================================================================

FROM alpine:latest

ENV DATABASE_URL=postgres://localhost/core?sslmode=disable
ENV ES_URL=http://localhost:9200

WORKDIR /root

COPY --from=build /go/src/github.com/astroband/astrologer/astrologer .
RUN chmod +x ./astrologer

CMD ["/root/astrologer", "ingest", "--", "50"]