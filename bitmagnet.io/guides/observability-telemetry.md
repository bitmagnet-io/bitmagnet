---
title: Observability & Telemetry
description: Observability & Telemetry features in bitmagnet
parent: Guides
layout: default
nav_order: 9
redirect_from:
  - /internals-development/observability-telemetry.html
---

# Observability & Telemetry

## Grafana stack & Prometheus integration

**bitmagnet** can integrate with the [Grafana stack](https://grafana.com/) and [Prometheus](https://prometheus.io/) for monitoring and building observability dashboards for the DHT crawler and other components. See the "Optional observability services" section of the [example docker compose configuration](https://github.com/bitmagnet-io/bitmagnet/blob/main/docker-compose.yml) and [example Grafana / Prometheus configuration files and a provisioned Grafana dashboard](https://github.com/bitmagnet-io/bitmagnet/tree/main/observability).

![Grafana dashboard](/assets/images/grafana-1.png)

The example integration includes:

- [Grafana](https://grafana.com/oss/grafana/) - A dashboarding and visualization tool
- [Grafana Agent](https://grafana.com/oss/agent/) - Collects metrics and logs, and forwards them to storage backends
- [Prometheus](https://prometheus.io/) - A time series database for metrics
- [Loki](https://grafana.com/oss/loki/) - A log aggregation system
- [Pyroscope](https://pyroscope.io/) - A continuous profiling tool
- [Postgres exporter](https://github.com/prometheus-community/postgres_exporter) - Exposes Postgres metrics to Prometheus

# Profiling with pprof

**bitmagnet** exposes [Go pprof](https://golang.org/pkg/net/http/pprof/) profiling endpoints at `/debug/pprof/*`, for example:

```sh
go tool pprof http://localhost:3333/debug/pprof/heap
```
