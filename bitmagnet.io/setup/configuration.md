---
title: Configuration
parent: Setup
layout: default
nav_order: 2
---

# Configuration

**bitmagnet** exposes quite a few configuration options. You shouldn't have to worry about most of them, but they are available for tinkering. If you're using the [example docker-compose file]({% link setup/installation.md %}#docker) then things _should_ "just work". I'll only cover only some of the more important options here:

- `postgres.host`, `postgres.name` `postgres.user` `postgres.password` (default: `localhost`, `bitmagnet`, `postgres`, _empty_): Set these values to configure connection to your Postgres database.
- `redis.addr`, `redis.db`, `redis.username`, `redis.password` (default: `localhost:6379`, `0`, _empty_, _empty_): Configure access to your Redis instance.
- `tmdb.api_key`: This is quite an important one, please [see below](#obtaining-a-tmdb-api-key) for more details.
- `dht_crawler.save_files` (default: `true`): If true, file metadata from the DHT crawler will be saved to the database. This provides more rich information about a torrent, but will use more disk space. If disk space is at a premium you may want to consider disabling this.
- `dht_crawler.save_files_threshold` (default: `50`): This parameter provides a compromise over disabling the saving of files altogether. Some torrents contain many thousands of files, which impacts performance and uses a lot of database disk space. When greater than `0`, this parameter will discard the files info when the number of files is greater than the threshold.
- `dht_crawler.save_pieces` (default: `false`): If true, the DHT crawler will save the pieces bytes from the torrent metadata. The pieces take up quite a lot of space, and aren't currently very useful, but they may be used by future features.
- `log.level` (default: `info`): If you're developing or just curious then you may want to set this to `debug`; note that `debug` output will be very verbose.
- `log.development` (default: `false`): If you're developing you may want to enable this flag to enable more verbose output such as stack traces.
- `log.json` (default: `false`): By default logs are output in a pretty format with colors; enable this flag if you'd prefer plain JSON.
- There are various configurable rate and concurrency limits associated with the DHT crawler. I won't detail them all here but you can see a list of them using the command below. I haven't spent a lot of time experimenting with these; they do have a big impact on the throughput of the crawler, and are available for tweaking. Some configurations of these values might result in issues such as too many open connections - I will aim to keep the defaults sensible, and am open to input on what those sensible defaults should be! The current default settings give me a throughput of around 5,000 discovered torrents per hour, though your mileage may vary based on numerous factors such as your hardware and internet connection.

To see a full list of available configuration options using the CLI, run:

```sh
bitmagnet config show
```

{% include callout_cli.md %}

For each configuration parameter available, this command will show:

- The path of the config key
- The Go type of the config key
- The currently configured value
- The default value
- Where the currently configured value has been sourced from (e.g. `default`, `./config.yml`, `env`)

## Specifying configuration values

Configuration paths are delimited by dots. If you're specifying configuration in a YAML file then each dot represents a nesting level, for example to configure `log.json`, `tmdb.api_key` and `http_server.cors.allowed_origins`:

```yaml
log:
  json: true
tmdb:
  api_key: my-api-key
http_server:
  cors:
    allowed_origins:
      - https://example1.com
      - https://example2.com
```

{: .note }
This is not a suggested configuration file, it's just an example of how to specify configuration values.

To configure these same values with environment variables, upper-case the path and replace all dots with underscores, for example:

```sh
LOG_JSON=true \
TMDB_API_KEY=my-api-key \
HTTP_SERVER_CORS_ALLOWED_ORIGINS=https://example1.com,https://example2.com \
  bitmagnet config show
```

## Configuration precedence

In order of precedence, configuration values will be read from:

- Environment variables
- `config.yml` in the current working directory
- `config.yml` in the [XDG-compliant](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) config location for the current user (for example on MacOS this is `~/Library/Application Support/bitmagnet/config.yml`)
- Default values

{: .warning }
Environment variables can be used to configure simple scalar types (strings, numbers, booleans) and slice types (arrays). For more complex configuration types such as maps you'll have to use YAML configuration. **bitmagnet** will exit with an error if it's unable to parse a provided configuration value.

## Obtaining a TMDB API Key

{: .highlight }
**bitmagnet** uses [the TMDB API](https://developer.themoviedb.org/docs){:target="\_blank"} to fetch metadata for movies and TV shows. By default you'll be sharing an API key with other users. If you're using this app and its content classifier heavily then you'll need to get a personal TMDB API key. Until you do this you'll see a warning message in the logs on startup, and you'll be limited to 1 TMDB API request per second. This is just about enough for running the DHT crawler, but if you're importing and classifying a lot of content this will be a major bottleneck. If many people are using this app with the default API key then that could add up to many requests per second, so please get your own API key if you are using this app more than casually!

Obtaining an API key is free and relatively easy, but you'll have to register for a TMDB account, provide them with some personal information such as contact details, a website URL (such as your GitHub account or social media profile) and a short description of your use case (**tip:** this app provides _"A content classifier that identifies movies and TV shows based on filenames"_). Once you've filled in the request form, approval should be instant.

[Synology have provided a full tutorial on obtaining a TMDB API key](https://kb.synology.com/en-au/DSM/tutorial/How_to_apply_for_a_personal_API_key_to_get_video_info){:target="\_blank"}.

Once you've obtained your API key you'll need to configure the `tmdb.api_key` value. Your rate limit will then default to 20 requests per second, which is well within [TMDB's stated fair usage limit](https://developer.themoviedb.org/docs/rate-limiting){:target="\_blank"}.
