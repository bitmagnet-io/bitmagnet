---
title: Tech Stack
description: The technology stack used in bitmagnet
parent: Guides
layout: default
nav_order: 10
redirect_from:
  - /internals-development.html
  - /internals-development/dht-crawler.html
  - /internals-development/tech-stack.html
---

# The Tech Stack

{: .highlight }
Are you an experienced developer with knowledge of GoLang, Postgres, TypeScript/Angular and/or BitTorrent protocols? I'm currently a lone developer with a full time job and many other commitments, and have been working on this in spare moments for the past few months. This project is too big for one person! If you're interested in contributing please [review the open issues](https://github.com/bitmagnet-io/bitmagnet/issues) and feel free to open a PR!

## Postgres

Postgres is the primary data store, and powers the search engine and message queue. These make use of several Postgres-specific features and extensions; as such, supporting other storage engines is likely to be complicated and is not a priority at the moment.

## GoLang Backend

Some key libraries used include:

- [anacrolix/torrent](https://github.com/anacrolix/torrent) not heavily used right now, but contains many useful BitTorrent utilities and could drive future features such as in-place seeding
- [fx](https://uber-go.github.io/fx/) for dependency injection and management of the application lifecycle
- [gin](https://gin-gonic.com/) for the HTTP server
- [goose](https://pressly.github.io/goose/) for database migrations
- [gorm](https://gorm.io/) for database access
- [gqlgen](https://gqlgen.com/) for the GraphQL server implementation
- [rex](https://github.com/hedhyw/rex) a regular expression library that makes some of the monstrous classification regexes more manageable
- [urfave/cli](https://cli.urfave.org/) for the command line interface
- [zap](https://github.com/uber-go/zap) for logging

## TypeScript/Angular Web UI

Using [Angular Material components](https://material.angular.io/). The web UI is embedded in the GoLang binary and served by the Gin web framework, and hence the build artifacts are committed into the repository.

## Nix dev shell

The repository includes a Nix shell for a reproducible development environment. To use the shell, simply [install Nix](https://nixos.org/download/) then run `nix develop` (or better still, use [nix-direnv](https://github.com/nix-community/nix-direnv) to use the included shell automatically.

## Other tooling

- The repository includes a [Taskfile](https://taskfile.dev/) containing several useful development scripts
- GitHub actions are used for CI, building the Docker image and for building this website

## Architecture & Lifecycle of the DHT Crawler

The DHT and BitTorrent protocols are (rather impenetrably) documented at [bittorrent.org](http://bittorrent.org/beps/bep_0000.html). Relevant resources include:

- [BEP 5: DHT Protocol](http://bittorrent.org/beps/bep_0005.html)
- [BEP 51: Infohash Indexing](https://www.bittorrent.org/beps/bep_0051.html)
- [BEP 33: DHT Scrapes](https://www.bittorrent.org/beps/bep_0033.html)
- [BEP 10: Extension Protocol](https://www.bittorrent.org/beps/bep_0010.html)
- [The Kademlia paper](https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf)

The rest of what I've figured out about how to implement a DHT crawler was cobbled together from [the now archived **magnetico** project](https://github.com/boramalper/magnetico) and [anacrolix's BitTorrent libraries](https://github.com/anacrolix).
