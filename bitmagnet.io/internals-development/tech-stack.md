---
title: Tech Stack
parent: Internals & Development
layout: default
nav_order: 1
---

# The Tech Stack

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

## Other tooling

- The repository includes a [Taskfile](https://taskfile.dev/) containing several useful development scripts
- GitHub actions are used for CI, building the Docker image and for building this website
