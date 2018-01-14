FROM bjornmagnusson/rpi-golang AS builder
COPY . $GOPATH/src/app
RUN go build -v led.go && \
    go test -v && \
    mv led /usr/local/bin && \
    chmod +x /usr/local/bin/led

FROM resin/rpi-raspbian:jessie-20170111
ENV PUSHOVER_TOKEN ""
ENV PUSHOVER_USER ""
COPY --from=builder /usr/local/bin/led .
EXPOSE 8080
ENTRYPOINT ["./led"]
