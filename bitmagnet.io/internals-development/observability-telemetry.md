---
title: Observability & Telemetry
parent: Internals & Development
layout: default
nav_order: 3
---

# Observability & Telemetry

## Grafana stack & Prometheus integration

**bitmagnet** can integrate with the [Grafana stack](https://grafana.com/){:target="\_blank"} and [Prometheus](https://prometheus.io/){:target="\_blank"} for monitoring and building observability dashboards for the DHT crawler and other components. See the "Optional observability services" section of the [example docker compose configuration](https://github.com/bitmagnet-io/bitmagnet/blob/main/docker-compose.yml){:target="\_blank"} and [example Grafana / Prometheus configuration files and a provisioned Grafana dashboard](https://github.com/bitmagnet-io/bitmagnet/tree/main/observability){:target="\_blank"}.

![Grafana dashboard](/assets/images/grafana-1.png)

The example integration includes:

- [Grafana](https://grafana.com/oss/grafana/){:target="\_blank"} - A dashboarding and visualization tool
- [Grafana Agent](https://grafana.com/oss/agent/){:target="\_blank"} - Collects metrics and logs, and forwards them to storage backends
- [Prometheus](https://prometheus.io/){:target="\_blank"} - A time series database for metrics
- [Loki](https://grafana.com/oss/loki/){:target="\_blank"} - A log aggregation system
- [Pyroscope](https://pyroscope.io/){:target="\_blank"} - A continuous profiling tool
- [Postgres exporter](https://github.com/prometheus-community/postgres_exporter){:target="\_blank"} - Exposes Postgres metrics to Prometheus

# Profiling with pprof

**bitmagnet** exposes [Go pprof](https://golang.org/pkg/net/http/pprof/){:target="\_blank"} profiling endpoints at `/debug/pprof/*`, for example:

```sh
go tool pprof http://localhost:3333/debug/pprof/heap
```
