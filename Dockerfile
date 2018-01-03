FROM resin/rpi-raspbian:jessie-20171227 as build

# Install Go
ENV GOVERSION 1.9.2
ENV GOTAR go$GOVERSION.linux-armv6l.tar.gz
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN apt-get update && \
    apt-get install -y --no-install-recommends wget && \
    wget https://storage.googleapis.com/golang/$GOTAR && \
    tar -C /usr/local -xzf $GOTAR

# Build application
RUN mkdir -p $GOPATH/src/app
WORKDIR $GOPATH/src/app
COPY . $GOPATH/src/app
RUN go build -v led.go && \
    mv led /usr/local/bin && \
    chmod +x /usr/local/bin/led && \
    go test

FROM resin/rpi-raspbian:jessie-20170111
ENV PUSHOVER_TOKEN ""
ENV PUSHOVER_USER ""
COPY --from=build /usr/local/bin/led .
EXPOSE 8080
ENTRYPOINT ["./led"]
