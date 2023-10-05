---
title: Installation
parent: Setup
layout: default
nav_order: 1
---

# Installation

## Docker

The quickest way to get up-and-running with **bitmagnet** is by using Docker with a docker-compose file. Adapt the following `docker-compose.yml` file to your needs:

```yml
services:
  bitmagnet:
    image: ghcr.io/bitmagnet-io/bitmagnet:latest
    container_name: bitmagnet
    ports:
      - "3333:3333"
    restart: unless-stopped
    environment:
      - POSTGRES_HOST=postgres
      - POSTGRES_PASSWORD=postgres
      - REDIS_ADDR=redis:6379
    #      - TMDB_API_KEY=your_api_key
    command:
      - worker
      - run
      - --keys=http_server
      - --keys=queue_server
      - --keys=dht_crawler # disable this line to run without DHT crawler
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

  postgres:
    image: postgres:16-alpine
    container_name: bitmagnet-postgres
    volumes:
      - ./data/postgres:/var/lib/postgresql/data
    #    ports:
    #      - "5432:5432" Expose this port if you'd like to dig around in the database
    restart: unless-stopped
    environment:
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=bitmagnet
      - PGUSER=postgres
    healthcheck:
      test:
        - CMD-SHELL
        - pg_isready
      start_period: 20s
      interval: 10s

  redis:
    image: redis:7-alpine
    container_name: bitmagnet-redis
    volumes:
      - ./data/redis:/data
    restart: unless-stopped
    entrypoint:
      - redis-server
      - --save 60 1
    healthcheck:
      test:
        - CMD
        - redis-cli
        - ping
      start_period: 20s
      interval: 10s
```

After running `docker-compose up -d` you should be able to access the web interface at `http://localhost:3333`. The DHT crawler should have started and you should see items appear in the web UI; try enabling auto-refresh to see live crawling and classification results.

## go install

You can also install **bitmagnet** natively with `go install github.com/bitmagnet-io/bitmagnet`. If you choose this method you will need to [configure]({% link setup/configuration.md %}) (at a minimum) Postgres and Redis instances for bitmagnet to connect to.

## Running the CLI

The **bitmagnet** CLI is the entrypoint into the application. Take note of the command needed to run the CLI, depending on your installation method.

- If you are using the docker-compose example above, you can run the CLI (while the stack is started) with `docker exec -it bitmagnet /bitmagnet`.
- If you installed bitmagnet with `go install`, you can run the CLI with `bitmagnet`.

When referring to CLI commands in the rest of the documentation, for simplicity we will use `bitmagnet`; please substitute this for the correct command. For example, to show the CLI help, run:

```sh
bitmagnet --help
```
