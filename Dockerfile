FROM resin/rpi-raspbian:jessie-20170111 as build

# Install Go
ENV GOVERSION 1.7.3
ENV GOTAR go$GOVERSION.linux-armv6l.tar.gz
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN apt-get update && \
    apt-get install -y --no-install-recommends wget && \
    wget https://storage.googleapis.com/golang/$GOTAR && \
    tar -C /usr/local -xzf $GOTAR && \
    rm -f $GOTAR && \
    # export PATH=$PATH:/usr/local/go/bin && \
    # export GOPATH=$PWD && \
    # go get -v github.com/rs/cors && \
    # go get -v github.com/stianeikeland/go-rpio && \
    # go get -v github.com/kidoman/embd && \
    apt-get remove -y --purge wget $(apt-mark showauto) && rm -rf /var/lib/apt/lists/*

# Build application
RUN mkdir -p $GOPATH/src/app
WORKDIR $GOPATH/src/app
COPY . $GOPATH/src/app
RUN go build -v led.go && \
    mv led /usr/local/bin && \
    chmod +x /usr/local/bin/led && \
    rm -rf /usr/local/go

FROM resin/rpi-raspbian:jessie-20170111
COPY --from=build /usr/local/bin/led .
EXPOSE 8080
ENTRYPOINT ["./led"]
