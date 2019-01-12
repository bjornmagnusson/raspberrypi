FROM golang:1.11.2 AS vendor
RUN curl https://glide.sh/get | sh && \
    glide init --non-interactive
ADD glide.yaml .
RUN glide up

FROM bjornmagnusson/rpi-golang AS builder
COPY --from=vendor /go/vendor vendor
ADD led.go .
RUN go build -v led.go

FROM bjornmagnusson/rpi-golang AS test
COPY --from=vendor /go/vendor vendor
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
