FROM quay.io/prometheus/busybox:latest
LABEL maintainer="theo.brigitte@gmail.com"

COPY example/simple.yml /etc/dummy_exporter.yml
COPY ./prometheus-dummy-exporter /bin/dummy_exporter

EXPOSE 9510
ENTRYPOINT ["/bin/dummy_exporter"]
CMD ["--config=/etc/dummy_exporter.yml"]
