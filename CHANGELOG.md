# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.11.0] - UNRELEASED

### Todo

- Docs: plugins.
- Docs: workers.
- Check web UI translations.
- Check errors.

### Added

- This CHANGELOG file (#?).
- A devcontainer configuration for local development, tested in VSCode and GoLand (#?).
- Implemented a plugin system, and port all core components to plugins (#?).
- Implemented Torznab profiles to allow for multiple Torznab configurations (#409).
- Catalan i18n (#404).
- A XXX category for Torznab queries (#432).
- Linked to Proxmox community script on external resources page (#430).
- Matching for SATRip and IPTVRip (#419).

### Fixed

- A bug whereby search queries prefixed with special characters such as "&" would result in an invalid tsquery and a database error (#423).

### Changed

- Upgraded go to 1.24.4.
- Refactored database persistence to use a single channel and API for all persistence operations (#?).
- Refactored torrent indexer to handle more torrents at a time while reducing bottlenecks (#422).
- Refactored DHT crawler to run an initial classification before persisting torrents, allowing crawl output to be available in the app more quickly (#?).
- Refactored info hash blocker to use Postgres streaming large objects for saving the bloom filter of blocked info hashes, which is significantly more memory efficient than saving these as Gorm models (#396).
- Deprecated the `bitmagnet worker run` command in favour of `bitmagnet start` (#?).
- Introduced stricter linting rules (#409, #421).

## [0.10.0] - 2025-03-02
