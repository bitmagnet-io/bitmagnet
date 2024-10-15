---
title: Installation
parent: Setup
layout: default
nav_order: 1
---

# Installation

## Docker

The quickest way to get up-and-running with **bitmagnet** is with [Docker Compose](https://docs.docker.com/compose/). The following `docker-compose.yml` is a minimal example. For a more full-featured example including VPN routing and observability services see the [docker compose configuration in the GitHub repository](https://github.com/bitmagnet-io/bitmagnet/blob/main/docker-compose.yml).

```yml
services:
  bitmagnet:
    image: ghcr.io/bitmagnet-io/bitmagnet:latest
    container_name: bitmagnet
    ports:
      # API and WebUI port:
      - "3333:3333"
      # BitTorrent ports:
      - "3334:3334/tcp"
      - "3334:3334/udp"
    restart: unless-stopped
    environment:
      - POSTGRES_HOST=postgres
      - POSTGRES_PASSWORD=postgres
    #      - TMDB_API_KEY=your_api_key
    volumes:
      - ./config:/root/.config/bitmagnet
    command:
      - worker
      - run
      - --keys=http_server
      - --keys=queue_server
      # disable the next line to run without DHT crawler
      - --keys=dht_crawler
    depends_on:
      postgres:
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
    shm_size: 1g
    healthcheck:
      test:
        - CMD-SHELL
        - pg_isready
      start_period: 20s
      interval: 10s
```

After running `docker compose up -d` you should be able to access the web interface at `http://localhost:3333`. The DHT crawler should have started and you should see items appear in the web UI within around a minute.

To upgrade your installation you can run:

```sh
docker compose down bitmagnet
docker pull ghcr.io/bitmagnet-io/bitmagnet:latest
docker compose up -d bitmagnet
```

## go install

You can also install **bitmagnet** natively with `go install github.com/bitmagnet-io/bitmagnet`. If you choose this method you will need to [configure]({% link setup/configuration.md %}) (at a minimum) a Postgres instance for bitmagnet to connect to.

## Running the CLI

The **bitmagnet** CLI is the entrypoint into the application. Take note of the command needed to run the CLI, depending on your installation method.

- If you are using the docker-compose example above, you can run the CLI (while the stack is started) with `docker exec -it bitmagnet bitmagnet`.
- If you installed bitmagnet with `go install`, you can run the CLI with `bitmagnet`.

When referring to CLI commands in the rest of the documentation, for simplicity we will use `bitmagnet`; please substitute this for the correct command. For example, to show the CLI help, run:

```sh
bitmagnet --help
```

## Starting **bitmagnet**

**bitmagnet** runs as multiple worker processes that can be started either individually or all at once. To start all workers, run:

```sh
bitmagnet worker run --all
```

Alternatively, specify individual workers to start:

```sh
bitmagnet worker run --keys=http_server,queue_server,dht_crawler
```
