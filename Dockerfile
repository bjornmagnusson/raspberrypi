FROM resin/rpi-raspbian

RUN wget https://storage.googleapis.com/golang/go1.7.3.linux-armv6l.tar.gz && \
    tar -C /usr/local -xzf go1.7.3.linux-armv6l.tar.gz && \
    export PATH=$PATH:/usr/local/go/bin && \    
    go get -v github.com/stianeikeland/go-rpio && \
    go get -v github.com/kidoman/embd && \
    go build -v led.go && \
    cp led /usr/local/bin && \
    chmod +x /usr/local/bin/led

EXPOSE 8080

CMD /usr/local/bin/led
