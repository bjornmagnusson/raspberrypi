FROM hypriot/rpi-golang

COPY led /usr/local/bin
RUN chmod +x /usr/local/bin/led

CMD /usr/local/bin/led
