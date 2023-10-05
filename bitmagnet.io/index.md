---
title: Home
layout: default
nav_order: -1
---

# bitmagnet

**A self-hosted BitTorrent indexer, DHT crawler, content classifier and torrent search engine with web UI, GraphQL API and Servarr stack integration.**

{: .warning-title }

> Important
>
> This software is currently in alpha. It is ready to preview some interesting and unique features, but there will likely be bugs, as well as API and database schema changes before the (currently theoretical) 1.0 release. If you'd like to support this project and help it gain momentum, **[please give it a star on GitHub](https://github.com/bitmagnet-io/bitmagnet){:target="\_blank"}**.
>
> [If you're interested in getting involved and you're a backend GoLang or frontend TypeScript/Angular developer, or you're knowledgeable about BitTorrent protocols then **I'd like to hear from you**](/internals-development.html) - let's get this thing over the line!

## DHT what now...?

The DHT crawler is **bitmagnet**'s killer feature that (currently) makes it unique. Well, almost unique, read on...

So what is it? You might be aware that you can enable DHT in your BitTorrent client, and that this allows you find peers who are announcing a torrent's hash to a Distributed Hash Table (DHT), rather than to a centralized tracker. DHT's lesser known feature is that it allows you to crawl the info hashes it knows about. This is how **bitmagnet**'s DHT crawler works works - it crawls the DHT network, requesting metadata about each info hash it discovers. It then further enriches this metadata by attempting to classify it and associate it with known pieces of content, such as movies and TV shows. It then allows you to search everything it has indexed.

This means that **bitmagnet** is not reliant on any external trackers or torrent indexers. It's a self-contained, self-hosted torrent indexer, connected via DHT to a global network of peers and constantly discovering new content.

The DHT crawler is _not quite_ unique to **bitmagnet**; another open-source project, [magnetico](https://github.com/boramalper/magnetico){:target="\_blank"} was first (as far as I know) to implement a usable DHT crawler, and was a crucial reference point for implementing this feature. However this project is no longer maintained, and does not provide the other features such as content classification, and integration with other software in the ecosystem, that greatly improve usability.

[You can find some more technical details about **bitmagnet**'s DHT crawler here](/internals-development/dht-crawler.html).

## Features & Roadmap

### Currently implemented features

- [x] A DHT crawler
- [x] A generic BitTorrent indexer: **bitmagnet** can index torrents from any source, not only the DHT network - currently this is only possible via [the `/import` endpoint](/tutorials/importing.html); more user-friendly methods are in the pipeline, see high-priority features below
- [x] A content classifier that can currently identify movie and television content, along with key related attributes such as language, resolution, source (BluRay, webrip etc.) and enriches this with data from [The Movie Database](https://www.themoviedb.org/)
- [x] [An import facility for ingesting torrents from any source, for example the RARBG backup](/tutorials/importing.html)
- [x] A torrent search engine
- [x] A GraphQL API: currently this provides a single search query; there is also an embedded GraphQL playground at `/graphql`
- [x] A web user interface implemented in Angular: currently this is a simple single-page application providing a user interface for search queries via the GraphQL API
- [x] [A Torznab-compatible endpoint for integration with the Serverr stack](/tutorials/servarr-integration.html)

### High priority features not yet implemented

- [ ] Classifiers for other types of content; enrich current classifiers and weed out incorrect classifications.
- [ ] Ordering of search results: the current alpha preview has no facility for specifying the ordering of results.
- [ ] Search performance optimisations: search is currently fast enough to be usable; it becomes more sluggish once millions of torrents have been indexed - there are some low-hanging fruit in terms of optimisation that will be a near-term priority.
- [ ] A monitoring API and WebUI dashboard showing things like crawler throughput, task queue, database size etc.
- [ ] Authentication, API keys, access levels etc.
- [ ] An admin API, and in general a more complete GraphQL API
- [ ] A more complete web UI
- [ ] Saved searches for content of particular interest, enabling custom feeds in addition to the following feature
- [ ] Smart deletion: there's a lot of crap out there; crawling DHT can quickly use lots of database disk space, and search becomes slower with millions of indexed torrents of which 90% are of no interest. A smart deletion feature would use saved searches to identify content that you're _not_ interested in, including but not limited to <abbr title="child sexual abuse material">CSAM</abbr>, and low quality content (such as low resolution movies). It would automatically delete associated metadata and add the info hash to a bloom filter, preventing the torrent from being re-indexed in future.
- [ ] Bi-directional integration with the [Prowlarr indexer proxy](https://prowlarr.com/){:target="\_blank"}: Currently **bitmagnet** can be added as an indexer in Prowlarr; bi-directional integration would allow **bitmagnet** to crawl content from any indexer configured in Prowlarr, unlocking many new sources of content
- [ ] More documentation and more tests!

### Pipe dream features

This is where things start to get a bit nebulous. For now all focus is on delivering the core features above, but some of these ideas could be explored in future:

- [ ] In-place seeding: identify files on your computer that are part of an indexed torrent, and allow them to be seeded in place after having moved, renamed or deleted parts of the torrent
- [ ] Integration with popular BitTorrent clients
- [ ] Federation of some sort: allow friends to connect instances and pool the indexing effort, perhaps involving crowd sourcing manual content curation to supplement the automated classifiers
- [ ] Something that looks like a decentralized private tracker; by this I probably mean something that's based partly on personal trust and manually weeding out any bad actors; I'd be wary of creating something that looks a bit like [Tribler](https://github.com/Tribler/tribler){:target="\_blank"}, which while an interesting project seems to have demonstrated that implementing trust, reputation and privacy at the protocol level carries too much overhead to be a compelling alternative to plain old BitTorrent, for all its imperfections
- [ ] Support for the [BitTorrent v2 protocol](https://blog.libtorrent.org/2020/09/bittorrent-v2/){:target="\_blank"}: It remains to be seen if wider adoption will ever make this a valuable feature
