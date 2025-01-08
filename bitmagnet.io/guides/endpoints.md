---
title: Endpoints
parent: Guides
layout: default
nav_order: 1
redirect_from:
  - /endpoints.html
---

# **bitmagnet** Endpoints

**bitmagnet** exposes functionality on a number of endpoints:

- `/` - Redirects to `/webui`
- `/webui` - Main web user interface
- `/graphql` - GraphQL API including the GraphiQL browser interface
- `/torznab/*` - Torznab API for integration compatible applications
- `/import` - Import API for adding new content to the library (see [the importing guide](/guides/import.html))
- `/metrics` - Prometheus metrics (see [the observability guide](/guides/observability-telemetry.html))
- `/debug/pprof/*` - Go pprof profiling endpoints (see [the observability guide](/guides/observability-telemetry.html))
- `/status` - Health check/status endpoint
