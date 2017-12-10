FROM golang:1.9.2

WORKDIR /go/src/usps

RUN go get -u github.com/golang/dep/cmd/dep

COPY . .

RUN go-wrapper download
