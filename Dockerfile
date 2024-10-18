FROM quay.io/prometheus/busybox:latest
LABEL maintainer="theo.brigitte@gmail.com"

COPY example/simple.yml /etc/prometheus-dummy-exporter.yaml
COPY ./prometheus-dummy-exporter /bin/prometheus-dummy-exporter

EXPOSE 9510
ENTRYPOINT ["/bin/prometheus-dummy-exporter"]
CMD ["--config=/etc/prometheus-dummy-exporter.yaml"]
