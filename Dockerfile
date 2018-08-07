FROM golang:alpine

RUN apk update && \
    apk add --no-cache git
    
RUN go get -v github.com/2Bot/2Bot2Go

WORKDIR /go/src/github.com/2Bot/2Bot2Go

#COPY . .

ENV GOBIN=/go

RUN go get -v ./... && go install -v ./...

VOLUME ["/go/emoji", "/go/config"]

WORKDIR /go

CMD ["/go/2Bot2Go", "-c"]
