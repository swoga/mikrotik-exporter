FROM alpine
LABEL org.opencontainers.image.source https://github.com/swoga/mikrotik-exporter

RUN apk add --no-cache tzdata

COPY mikrotik-exporter /bin/mikrotik-exporter
COPY example.yml /etc/mikrotik-exporter/config.yml
COPY dist/modules /etc/mikrotik-exporter/modules
COPY dist/modules /modules

EXPOSE 9436

ENTRYPOINT ["/bin/mikrotik-exporter"]
CMD ["--config.file=/etc/mikrotik-exporter/config.yml"]
