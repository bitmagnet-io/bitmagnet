---
title: MatchLab (optional)
description: Optional external IMDb/episode matching with MatchLab
parent: Setup
layout: default
nav_order: 3
---

# Optional MatchLab integration

[MatchLab](https://github.com/Dizzybog/MatchLab) is a separate, self-hosted torrent parsing and IMDb/AniDB matching API. Bitmagnet can use MatchLab to identify video torrents without replacing its DHT crawler, PostgreSQL storage, existing classifier/queues or web interface.

**The integration is disabled by default.** An existing Bitmagnet installation continues to use native content lookup and TMDB title matching when no MatchLab URL is configured.

## Configure and run

Deploy MatchLab separately and make sure its API base URL is reachable **from the Bitmagnet process/container**, usually over a private Docker network. The API must respond successfully to `GET /health/ready`.

Add the following environment variables to the existing Bitmagnet service:

```yaml
services:
  bitmagnet:
    environment:
      MATCHLAB_URL: "http://matchlab:8765"
      MATCHLAB_TIMEOUT_SECONDS: "8"
      MATCHLAB_MAX_CONCURRENT: "4"
```

Preserve the existing Bitmagnet image, database, ports, volumes, network and other settings. MatchLab is not bundled with Bitmagnet. Set `MATCHLAB_URL` to the **base URL**, not the `/parse` or `/match` path. Keep the service on a trusted network or protect externally accessible API endpoints.

### Configuration keys

- `matchlab.url` / `MATCHLAB_URL`: MatchLab HTTP(S) base URL (empty by default; disables the integration).
- `matchlab.timeout_seconds` / `MATCHLAB_TIMEOUT_SECONDS`: Per-request timeout, default `8` seconds (allowed 1–120).
- `matchlab.max_concurrent` / `MATCHLAB_MAX_CONCURRENT`: Maximum concurrent MatchLab calls, default `4` (allowed 1–64).
- `tmdb.enabled` / `TMDB_ENABLED`: Enable TMDB metadata enrichment, default `true`.
- `classifier.concurrency` / `CLASSIFIER_CONCURRENCY`: Overall classifier concurrency, default `10`.

The same settings may be supplied through Bitmagnet's normal `config.yml` configuration:

```yaml
matchlab:
  url: http://matchlab:8765
  timeout_seconds: 8
  max_concurrent: 4
tmdb:
  enabled: true
```

For general configuration precedence and database settings, see [Configuration]({% link setup/configuration.md %}).

## Classification behavior

With a MatchLab URL set:

1. Bitmagnet uses its normal classifier to inspect torrent files and determine the content category. Non-video classification remains unchanged.
2. Movie/TV identity matching uses the **original torrent name** and MatchLab's `/match` endpoint. Verified IMDb IDs and season/episode coordinates come from this result, not from local title guesses.
3. MatchLab `/parse` can supplement the display title and release year; it cannot override verified episode coordinates.
4. When MatchLab verifies an IMDb ID, TMDB's **find by external ID** and details endpoints may enrich the record with artwork, descriptions and other metadata. TMDB title search does not serve as a fallback matcher in this mode.
5. An IMDb identity without TMDB metadata is still retained as a minimal IMDb-keyed content record when enrichment cannot identify the ID.
6. No verified IMDb identity is fabricated when MatchLab returns no match. AniDB-only coordinates may still indicate TV/series content.

Existing movie/TV hints and previously attached identities are revalidated against MatchLab when they pass through classification with MatchLab configured. Previously classified torrents are **not** automatically submitted en masse when switching on the integration; use Bitmagnet's existing queue reprocess/rematch action to revisit historical records.

### Failure handling

- **Verified no match:** The content remains unmatched. There is no fallback title matcher while MatchLab is configured.
- **Provider timeout, HTTP error, or malformed response:** The classifier reports an error, allowing the normal queue retry/backoff path to handle the job.
- **Parsing metadata unavailable after a verified match:** Preserve the verified identity and use available title/year metadata.
- **TMDB metadata unavailable for a verified IMDb ID:** Preserve a minimal IMDb-keyed record when the ID is not found; retry transient TMDB transport/server failures.
- **Concurrent jobs:** Requests share the configured outbound concurrency limit; configure this alongside the MatchLab server capacity.

## Dashboard and operations

When enabled, MatchLab appears in the dashboard's **Health** page, alongside the existing TMDB and database status. Its readiness check uses MatchLab's `/health/ready` endpoint. Use the existing **Queues** interface to inspect retries/failures and reprocess existing content.

No extra Bitmagnet database schema migration or custom MatchLab dashboard is needed for this integration.

## Disable MatchLab / return to native matching

1. Remove `MATCHLAB_URL` and any `matchlab.url` value from configuration.
2. Restart Bitmagnet. New queue work will use the original native content lookup and TMDB title matching.
3. Already saved IMDb-keyed content remains in PostgreSQL and can be looked up by its canonical or alternate IMDb identifier.
4. To **reidentify historical torrents with the native matcher**, use the existing Queues reprocess/rematch action. Bitmagnet deliberately does not automatically reprocess millions of old rows at startup.

Jobs queued before disconnection can contain a stale `matchlab_enabled` flag. When no MatchLab client is configured, the classifier ignores that stale flag and respects the normal `apis_enabled` and `local_search_enabled` switches. There is no MatchLab network call or health dependency in this mode.

## Tests and compatibility

The integration has tests for:

- API validation, invalid URLs, malformed identities and provider failures.
- Verified movie/TV matching and episode metadata.
- IMDb-to-TMDB metadata enrichment without title-search fallback.
- Retrying failed classification jobs and honoring cancellation.
- MatchLab-to-native recovery even with previously queued jobs.
- Retaining previously stored IMDb-keyed content.
- Dashboard Health integration.

See [the MatchLab project](https://github.com/Dizzybog/MatchLab) for installation and API implementation. For normal Bitmagnet operation see the [main project documentation](https://bitmagnet.io/).
