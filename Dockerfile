FROM fedora:26

RUN dnf -y update --refresh && \
    dnf -y install golang git && \
    dnf clean all

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin
RUN go get -u -v github.com/kardianos/govendor

ADD . $GOPATH/src/github.com/vrutkovs/polna
WORKDIR $GOPATH/src/github.com/vrutkovs/polna

RUN govendor sync -v && \
    go build && \
    go install && \
    mkdir upload

ENV GIN_MODE release

EXPOSE 8080

CMD ["/go/bin/polna"]