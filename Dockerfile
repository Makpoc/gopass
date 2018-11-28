FROM golang:latest

ADD . /go/src/github.com/Makpoc/gopass
WORKDIR /go/src/github.com/Makpoc/gopass/ui/gopass-web/
RUN go build .

ENV PORT=8080
EXPOSE $PORT

ENTRYPOINT /go/src/github.com/Makpoc/gopass/ui/gopass-web/gopass-web
