<p align="center">
    <img src="assets/prometheus-dummy-exporter.png" alt="prometheus-dummy-exporter" height="100px">
</p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/TheoBrigitte/prometheus-dummy-exporter"><img src="https://goreportcard.com/badge/github.com/TheoBrigitte/prometheus-dummy-exporter" alt="Go Report Card"></a>
  <a href="https://github.com/TheoBrigitte/prometheus-dummy-exporter/releases"><img src="https://img.shields.io/github/release/TheoBrigitte/prometheus-dummy-exporter"></a>
  <a href="https://hub.docker.com/r/theo01/prometheus-dummy-exporter"><img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/theo01/prometheus-dummy-exporter"></a>
  <a href="https://github.com/TheoBrigitte/prometheus-dummy-exporter/actions/workflows/test.yaml"><img src="https://github.com/TheoBrigitte/prometheus-dummy-exporter/actions/workflows/test.yaml/badge.svg?branch=main" alt="Github action"></a>
</p>


## Overview

`prometheus-dummy-exporter` exports meaningless and configurable metrics for [Prometheus](https://prometheus.io/).
It can be used for performance testing or developement in the Prometheus ecosystem.

## Install

### Binary

Go to https://github.com/TheoBrigitte/prometheus-dummy-exporter/releases

Or use `go install`

```bash
go install github.com/TheoBrigitte/prometheus-dummy-exporter
```

### Building from source

```bash
git clone https://github.com/TheoBrigitte/prometheus-dummy-exporter
cd prometheus-dummy-exporter
make build
```

### Docker container

https://hub.docker.com/r/theo01/prometheus-dummy-exporter

```bash
docker run -p 9510:9510 -v /PATH/TO/config.yaml:/etc/prometheus-dummy-exporter.yaml theo01/prometheus-dummy-exporter
```

### Kubernetes

```bash
kubectl apply -f https://raw.githubusercontent.com/TheoBrigitte/prometheus-dummy-exporter/main/kubernetes/manifest.yaml
```

## Usage

```
$ ./prometheus-dummy-exporter --help
usage: prometheus-dummy-exporter [<flags>]


Flags:
  -h, --[no-]help         Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=":9510"
                          Address to listen on for web interface and telemetry
      --web.telemetry-path="/metrics"
                          Path under which to expose metrics.
      --config=""         Path to config file
      --log.level="info"  Log level: [panic fatal error warning info debug trace]
      --[no-]version      Show application version.

```

Configuration format is below.

```yaml
# config.yml
namespace: <string>
metrics:
- name: <string>
  # support types are "counter" and "gauge"
  type: <string>
  # number of metrics
  size: <integer>
  labels:
    # label maps, values are selected using round robin
    <string>: [<string>, ...]
```

Example

```yaml
namespace: simple
metrics:
- name: alice
  type: counter
  size: 2
- name: bob
  type: gauge
  size: 8
  labels:
    foo: [one, two]
    bar: [three, four]
```

```
$ curl -s localhost:9510/metrics | egrep ^simple
simple_alice{id="0"} 1
simple_alice{id="1"} 1
simple_bob{bar="four",foo="two",id="1"} 0.008361162894281814
simple_bob{bar="four",foo="two",id="3"} 0.2928618665869591
simple_bob{bar="four",foo="two",id="5"} 0.7006760266267518
simple_bob{bar="four",foo="two",id="7"} 0.33883308503299564
simple_bob{bar="three",foo="one",id="0"} 0.23768179587904828
simple_bob{bar="three",foo="one",id="2"} 0.4605768559959702
simple_bob{bar="three",foo="one",id="4"} 0.8750007790390423
simple_bob{bar="three",foo="one",id="6"} 0.1504911738984869
```

## License

MIT
