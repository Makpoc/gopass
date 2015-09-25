FROM golang:latest

ADD . /go/src/github.com/makpoc/gopass
WORKDIR /go/src/github.com/makpoc/gopass/ui/gopass-web/
RUN go build .

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT /go/src/github.com/makpoc/gopass/ui/gopass-web/gopass-web
