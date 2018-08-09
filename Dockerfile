FROM golang:alpine

RUN apk update && \
    apk add --no-cache git

RUN go get github.com/2Bot/2Bot2Go

RUN mkdir -p /go/emoji && mkdir -p /go/config

WORKDIR /go/src/github.com/2Bot/2Bot2Go

ENV GOBIN=/go/

RUN go install -v ./...

VOLUME ["/go/emoji", "/go/config/config.toml"]

WORKDIR /go

CMD ["/go/2Bot2Go", "-c"]
