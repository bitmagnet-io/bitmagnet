# bitmagnet

A self-hosted BitTorrent indexer, DHT crawler, content classifier and torrent search engine with web UI, GraphQL API and Servarr stack integration.

Visit the website at [bitmagnet.io](https://bitmagnet.io).

## Optional MatchLab torrent matching

Bitmagnet can optionally use [MatchLab](https://github.com/Dizzybog/MatchLab) as an external IMDb identity matcher for movies and TV releases. This integration uses the existing classifier and queue pipeline; it does not replace Bitmagnet's DHT crawler, PostgreSQL database, web UI, or API.

With `MATCHLAB_URL` configured:

- The classifier sends original torrent names to MatchLab and uses its verified IMDb ID and season/episode coordinates for movie/TV identification, instead of matching titles locally or searching TMDB by title.
- TMDB may enrich the verified IMDb identity with artwork and metadata via external-ID lookup; TMDB remains a metadata provider rather than a second title matcher.
- MatchLab appears in dashboard Health, and failed requests retain Bitmagnet's normal queue retry semantics.
- Already-classified torrents can be reprocessed explicitly from the existing Queues interface; existing data is not automatically rematched en masse on startup.

Without `MATCHLAB_URL`, the original local/TMDB classification workflow runs as before. Disabling MatchLab after previous use does not delete existing torrent/content records; an explicit reprocess/rematch is needed to reclassify historical rows using the native engine.

**Configuration:** `MATCHLAB_URL=http://matchlab:8765` (optional), `MATCHLAB_TIMEOUT_SECONDS=8` and `MATCHLAB_MAX_CONCURRENT=4`.

See the [MatchLab integration guide](./bitmagnet.io/setup/matchlab.md) for Docker/network setup, configuration, failure handling and migration back to native classification.
