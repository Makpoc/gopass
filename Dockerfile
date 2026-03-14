FROM golang:1.25

ADD . /go/src/github.com/Makpoc/gopass
WORKDIR /go/src/github.com/Makpoc/gopass/ui/gopass-web/
RUN go build .

ENV GOPASS_PORT=8080
EXPOSE $GOPASS_PORT

ENTRYPOINT /go/src/github.com/Makpoc/gopass/ui/gopass-web/gopass-web
