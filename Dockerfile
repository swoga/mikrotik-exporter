FROM alpine
LABEL org.opencontainers.image.source https://github.com/swoga/mikrotik-exporter

COPY mikrotik-exporter /bin/mikrotik-exporter
COPY example.yml /etc/mikrotik-exporter/config.yml

EXPOSE 9436

ENTRYPOINT ["/bin/mikrotik-exporter"]
CMD ["--config.file=/etc/mikrotik-exporter/config.yml"]