---
title: Backup, Restore & Merge
parent: Tutorials
layout: default
nav_order: 2
---

# Backup, Restore & Merge

It's a good idea to take periodic backups of the **bitmagnet** database. A system crash might corrupt the database, in which case you could see the dreaded `PANIC: could not locate a valid checkpoint record` when starting Postgres.

Perhaps you'd like to move your **bitmagnet** installation to a new server, or you'd like to merge the data from two **bitmagnet** installations.

This tutorial will show you how to backup, restore and merge **bitmagnet** databases.

{: .note-title }

> Pre-requisites
>
> - [x] You'll need to have `pg_dump` and `psql` installed. These are part of the PostgreSQL package. Use Google to find out how to install these tools on your operating system.

## Taking a backup

The following command will take a backup of the critical **bitmagnet** data and save it to a file named `export.sql`. (note this is not a full backup of the database which would include creation of tables, indexes etc.). By exporting with the `--data-only` flag the resulting file can be imported into a new or existing installation, after **bitmagnet** has run its migrations to set up the database and tables.

Please refer to [the `pg_dump` documentation](https://www.postgresql.org/docs/current/app-pgdump.html) and ensure to specify the correct values (e.g. `host`, `username` and `password`) for the source database.

```sh
pg_dump \
        --column-inserts \
        --data-only \
        --on-conflict-do-nothing \
        --rows-per-insert=1000 \
        --table=metadata_sources \
        --table=content \
        --table=content_attributes \
        --table=content_collections \
        --table=content_collections_content \
        --table=torrent_sources \
        --table=torrents \
        --table=torrent_files \
        --table=torrent_hints \
        --table=torrent_contents \
        --table=torrents_torrent_sources \
        --table=key_values \
        bitmagnet \
        > backup.sql
```

## Restoring a backup, or merging into another **bitmagnet** instance

First, ensure you have a target **bitmagnet** instance up and running, _of the same version from which the backup was taken_.

The following command will import the backup file into the target database, merging the data with any existing data.

Please refer to [the `psql` documentation](https://www.postgresql.org/docs/current/app-psql.html) and ensure to specify the correct values (e.g. `host`, `username` and `password`) for the target database.

```sh
psql bitmagnet < backup.sql
```
