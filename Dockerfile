FROM golang:1.14.0-buster

WORKDIR /go/src/app
COPY src/ .
COPY src/pages /go/src/pages
RUN go get -u github.com/gorilla/mux
