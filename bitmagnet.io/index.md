---
title: Home
layout: default
nav_order: -1
---

# bitmagnet

**A self-hosted BitTorrent indexer, DHT crawler, content classifier and torrent search engine with web UI, GraphQL API and Servarr stack integration.**

![WebUI screenshot](/assets/images/webui-1.png)

{: .warning-title }

> Important
>
> This software is currently in alpha. It is ready to preview some interesting and unique features, but there will likely be bugs, as well as API and database schema changes before the (currently theoretical) 1.0 release. If you'd like to support this project and help it gain momentum, please **[give it a star on GitHub](https://github.com/bitmagnet-io/bitmagnet)** or **[sponsor it on OpenCollective](https://opencollective.com/bitmagnet)**.

## DHT what now...?

The DHT crawler is **bitmagnet**'s killer feature that makes it unique. So what is it?

You might be aware that you can enable DHT in your BitTorrent client, and that this allows you find peers who are announcing a torrent's hash to a Distributed Hash Table (DHT), rather than to a centralized tracker. DHT's lesser known feature is that it allows you to crawl the info hashes it knows about. This is how **bitmagnet**'s DHT crawler works - it crawls the DHT network, requesting metadata about each info hash it discovers. It then further enriches this metadata by attempting to classify it and associate it with known pieces of content, such as movies and TV shows. It then allows you to search everything it has indexed.

This means that **bitmagnet** is not reliant on any external trackers or torrent indexers. It's a self-contained, self-hosted torrent indexer, connected via DHT to a global network of peers and constantly discovering new content.

## Features & Roadmap

### Currently implemented features

- [x] A DHT crawler and protocol implementation
- [x] A generic BitTorrent indexer: **bitmagnet** can index torrents from any source, not only the DHT network - currently this is only possible via [the `/import` endpoint](/guides/import.html); more user-friendly methods are in the pipeline, see high-priority features below
- [x] A highly customizable <a href="/guides/classifier.html">content classifier</a> that can currently identify many types of content, along with key related attributes such as language, resolution, source (BluRay, webrip etc.) and enriches this with data from sources including [The Movie Database](https://www.themoviedb.org/)
- [x] [An import facility for ingesting torrents from any source, for example the RARBG backup](/guides/import.html)
- [x] A torrent search engine
- [x] A GraphQL API: currently this provides a single search query; there is also an embedded GraphQL playground at `/graphql`
- [x] A web user interface implemented in Angular: currently this is a simple single-page application providing a user interface for search queries via the GraphQL API
- [x] [A Torznab-compatible endpoint for integration with the Serverr stack](/guides/servarr-integration.html)
- [x] A WebUI dashboard for monitoring and administration

### High priority features not yet implemented

- [ ] Authentication, API keys, access levels etc.
- [ ] An admin API, and in general a more complete GraphQL API
- [ ] Saved searches for content of particular interest, enabling custom feeds in addition to the following feature
- [ ] Bi-directional integration with the [Prowlarr indexer proxy](https://prowlarr.com/): Currently **bitmagnet** can be added as an indexer in Prowlarr; bi-directional integration would allow **bitmagnet** to crawl content from any indexer configured in Prowlarr, unlocking many new sources of content
- [ ] More documentation and more tests!

### Pipe dream features

This is where things start to get a bit nebulous. For now all focus is on delivering the core features above, but some of these ideas could be explored in future:

- [ ] In-place seeding: identify files on your computer that are part of an indexed torrent, and allow them to be seeded in place after having moved, renamed or deleted parts of the torrent
- [ ] Integration with popular BitTorrent clients
- [ ] Federation of some sort: allow friends to connect instances and pool the indexing effort, perhaps involving crowd sourcing manual content curation to supplement the automated classifiers
- [ ] Something that looks like a decentralized private tracker; by this I probably mean something that's based partly on personal trust and manually weeding out any bad actors; I'd be wary of creating something that looks a bit like [Tribler](https://github.com/Tribler/tribler), which while an interesting project seems to have demonstrated that implementing trust, reputation and privacy at the protocol level carries too much overhead to be a compelling alternative to plain old BitTorrent, for all its imperfections
- [ ] Support for the [BitTorrent v2 protocol](https://blog.libtorrent.org/2020/09/bittorrent-v2/): It remains to be seen if wider adoption will ever make this a valuable feature
