FROM resin/rpi-raspbian

ENV GOVERSION 1.7.3
RUN apt-get update && \
    apt-get install wget && \
    apt-get install git && \
    wget https://storage.googleapis.com/golang/go$GOVERSION.linux-armv6l.tar.gz && \
    tar -C /usr/local -xzf go$GOVERSION.linux-armv6l.tar.gz && \
    rm -rf /var/lib/apt/lists
RUN export PATH=$PATH:/usr/local/go/bin && \
    export GOPATH=$PWD && \
    go get -v github.com/rs/cors && \
    go get -v github.com/stianeikeland/go-rpio && \
    go get -v github.com/kidoman/embd
COPY led.go .
RUN export PATH=$PATH:/usr/local/go/bin && \
    export GOPATH=$PWD && \
    go build -v led.go && \
    cp led /usr/local/bin && \
    chmod +x /usr/local/bin/led && \
    rm -rf /usr/local/go && \
    rm -f go$GOVERSION.linux-armv6l.tar.gz

EXPOSE 8080

CMD ["/usr/local/bin/led"]
