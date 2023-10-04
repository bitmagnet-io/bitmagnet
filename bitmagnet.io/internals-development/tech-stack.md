---
title: Tech Stack
parent: Internals & Development
layout: default
nav_order: 1
---

# The Tech Stack

## Postgres

Postgres is the primary data store, and powers the search engine. The search engine makes use of several Postgres-specific features and extensions; as such, supporting other storage engines is likely to be complicated and is not a priority at the moment.

## Redis

Redis is currently used only for the task queue, which currently only handles classification jobs (though other types of jobs and schedules that would make use of this are in the high-priority pipeline). The [asynq](https://github.com/hibiken/asynq){:target="\_blank"} library is used for this. I deliberated over adding Redis as a dependency but it was the most pragmatic solution; I wasn't able to find a mature off-the-shelf Postgres-backed solution that works well with GoLang. On the one hand if there was such a solution then I'd consider removing this dependency; on the other hand I can see Redis being useful for other features, such as improving support for distributing workers across multiple nodes, and potentially moving the DHT staging and routing tables to Redis.

## GoLang Backend

Some key libraries used include:

- [anacrolix/torrent](github.com/anacrolix/torrent){:target="\_blank"} not heavily used right now, but contains many useful BitTorrent utilities and could drive future features such as in-place seeding
- [asynq](https://github.com/hibiken/asynq){:target="\_blank"} for the task queue
- [fx](https://uber-go.github.io/fx/){:target="\_blank"} for dependency injection and management of the application lifecycle
- [gin](https://gin-gonic.com/){:target="\_blank"} for the HTTP server
- [goose](https://pressly.github.io/goose/){:target="\_blank"} for database migrations
- [gorm](https://gorm.io/){:target="\_blank"} for database access
- [gqlgen](https://gqlgen.com/){:target="\_blank"} for the GraphQL server implementation
- [rex](https://github.com/hedhyw/rex){:target="\_blank"} a regular expression library that makes some of the monstrous classification regexes more manageable
- [urfave/cli](https://cli.urfave.org/){:target="\_blank"} for the command line interface
- [zap](https://github.com/uber-go/zap){:target="\_blank"} for logging

## TypeScript/Angular Web UI

Using [Angular Material components](https://material.angular.io/){:target="\_blank"}. The web UI is embedded in the GoLang binary and served by the Gin web framework, and hence the build artifacts are committed into the repository.

## Other tooling

- The repository includes a [Taskfile](https://taskfile.dev/){:target="\_blank"} containing several useful development scripts
- GitHub actions are used for CI, building the Docker image and for building this website
