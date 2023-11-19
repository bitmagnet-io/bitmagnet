---
title: Endpoints
layout: default
nav_order: 2
---

# **bitmagnet** Endpoints

**bitmagnet** exposes functionality on a number of endpoints:

- `/` - Main web user interface
- `/graphql` - GraphQL API including the GraphiQL browser interface
- `/torznab/*` - Torznab API for integration compatible applications
- `/import` - Import API for adding new content to the library (see [the importing tutorial](/tutorials/importing.html))
- `/metrics` - Prometheus metrics (see [the observability guide](/internals-development/observability-telemetry.html))
- `/debug/pprof/*` - Go pprof profiling endpoints (see [the observability guide](/internals-development/observability-telemetry.html))
- `/asynqmon` - [Web UI for the Asynq task queue](https://github.com/hibiken/asynqmon){:target="\_blank"}
- `/status` - Health check/status endpoint
