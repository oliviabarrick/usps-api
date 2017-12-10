FROM golang:1.9.2

WORKDIR /go/src/github.com/justinbarrick/usps-api

COPY . .

RUN go get -d -t ./...
