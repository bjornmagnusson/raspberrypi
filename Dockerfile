FROM resin/rpi-raspbian:jessie-20161228

COPY led /usr/local/bin
RUN chmod +x /usr/local/bin/led

EXPOSE 8080

CMD /usr/local/bin/led
