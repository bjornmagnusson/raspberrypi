FROM golang:1.11.2 AS vendor
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN mkdir -p /go/src/github.com/bjornmagnusson/raspberrypi
WORKDIR /go/src/github.com/bjornmagnusson/raspberrypi
ADD Gopkg.toml .
ADD led.go .
RUN dep ensure -v

FROM bjornmagnusson/rpi-golang AS builder
COPY --from=vendor /go/src/github.com/bjornmagnusson/raspberrypi/vendor vendor
ADD led.go .
RUN go build -v led.go

FROM bjornmagnusson/rpi-golang AS test
COPY --from=vendor /go/src/github.com/bjornmagnusson/raspberrypi/vendor vendor
ADD led.go .
ADD led_test.go .
RUN go test -v

FROM balenalib/rpi-raspbian:jessie-20181201 AS dist
ENV PUSHOVER_TOKEN ""
ENV PUSHOVER_USER ""
COPY --from=builder /go/src/app/led .
RUN chmod +x led
EXPOSE 8080
ENTRYPOINT ["./led"]
