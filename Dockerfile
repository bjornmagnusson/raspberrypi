FROM hypriot/rpi-golang

COPY led /usr/local/bin
RUN chmod +x /usr/local/bin/led

EXPOSE 8080

CMD /usr/local/bin/led
