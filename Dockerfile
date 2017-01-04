FROM hypriot/rpi-golang

COPY led .

CMD [./led]
