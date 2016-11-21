FROM golang:1.7.3
MAINTAINER bruce.downs@jivesoftware.com

RUN go get github.com/bruceadowns/miru-syslog

EXPOSE 514 8081

ENTRYPOINT /go/bin/miru-syslog
