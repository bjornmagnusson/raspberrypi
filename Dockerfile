FROM golang:1.12.5 AS builder
RUN mkdir -p /app
WORKDIR /app
ADD go.mod .
ADD go.sum .
ADD led.go .
ENV GOPROXY=https://gocenter.io
RUN go build -v

FROM golang:1.12.5 AS test
RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /go/pkg/mod /go/pkg/mod
COPY --from=builder /app/raspberrypi /usr/local/bin/
ADD go.mod .
ADD go.sum .
ADD led.go .
ADD led_test.go .
RUN go test -v && \
    raspberrypi -demo=true -num=3

FROM balenalib/rpi-raspbian:jessie-20181201 AS dist
ENV PUSHOVER_TOKEN ""
ENV PUSHOVER_USER ""
COPY --from=builder /app/raspberrypi /usr/local/bin/
RUN chmod +x /usr/local/bin/raspberrypi
EXPOSE 8080
ENTRYPOINT ["raspberrypi"]
