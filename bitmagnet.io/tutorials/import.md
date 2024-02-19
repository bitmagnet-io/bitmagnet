---
title: Import
parent: Tutorials
layout: default
nav_order: 3
redirect_from:
  - /tutorials/importing.html
---

# Import

{: .warning-title }

> Important
>
> Before continuing with this tutorial, please [obtain and configure a personal TMDB API key]({% link setup/configuration.md %}#obtaining-a-tmdb-api-key).

**bitmagnet** includes an import endpoint at `/import`; this can be used for importing Torrent files from any source.

{: .note }

> A proper schema is needed for this endpoint, along with improved input validation. There isn't currently a way to import a torrent along with information about the files it contains (which is optional in **bitmagnet**). If an imported torrent is later discovered by the DHT crawler then its associated file info would be saved at that point.

## Example: The RARBG backup

For the purposes of this tutorial we'll use the RARBG SQLite backup, but you can adapt this example to any suitable data source.

{: .note-title }

> Pre-requisites
>
> - [x] You have [obtained and configured a personal TMDB API key]({% link setup/configuration.md %}#obtaining-a-tmdb-api-key)
> - [x] You have obtained a copy of the RARBG SQLite backup (I can't assist you in getting a copy of this, but it's generally available)
> - [x] You have [installed the SQLite3 CLI](https://www.tutorialspoint.com/sqlite/sqlite_installation.htm){:target="\_blank"}
> - [x] You have [installed jq](https://jqlang.github.io/jq/download/){:target="\_blank"}

Let's start by write a SQLite query in a file named `rarbg-import.sql`. This will extract the data we need and get it looking a bit more like the format that **bitmagnet** expects. The following is a starting point, please adapt it to your requirements:

```sql
select
  hash as infoHash,
  title as name,
  size,
-- map the RARBG category to a valid bitmagnet content type:
  case
    when cat like 'ebooks%' then 'ebook'
    when cat like 'games%' then 'game'
    when cat like 'movies%' then 'movie'
    when cat like 'tv%' then 'tv_show'
    when cat like 'music%' then 'music'
    when cat like 'software%' then 'software'
    when cat = 'xxx' then 'xxx'
  end as contentType,
-- we can give the classifier an easier job if we already know some characteristics of the content:
  case
    when cat like '%_4k' then 'V2160p'
    when cat like '%_720' then 'V720p'
    when cat like '%_SD' then 'V480p'
  end as videoResolution,
  case
    when cat like '%_bd_%' then 'BluRay'
  end as videoSource,
  case
    when cat like '%_bd_full' then 'BRDISK'
    when cat like '%_bd_remux' then 'REMUX'
  end as videoModifier,
  case
    when cat like '%_x264%' then 'x264'
    when cat like '%_x265%' then 'x265'
    when cat like '%_xvid%' then 'XviD'
  end as videoCodec,
  case
    when cat like '%_3D' then 'V3D'
  end as video3D,
  imdb,
-- convert the dt field to a valid ISO date string:
  (substr(dt, 0, 11) || 'T' || substr(dt, 12) || '.000Z') as publishedAt
  from items
  where

-- the following lines are optional;
-- it's recommended to review which categories are of interest,
-- as filtering unwanted and low quality content at the import stage will improve the app experience
    cat not like '%_720' and
    cat not like '%_SD' and
    cat not like 'software%' and
    cat not like 'games%' and
-- I won't judge you if you disable the following line;
-- bear in mind there is a *lot* of this in the RARBG backup
    cat != 'xxx' and

--
    true
-- a random-ish but deterministic order reduces the chances of the resolver duplicating its work:
  order by infoHash

-- you may want to enable the following line while testing your query
-- limit 100
```

You can try running this query in your favourite database explorer, or using the SQLite3 CLI.

So far we've got the data looking almost like we need it. We now need to make a few final tweaks before piping it into **bitmagnet**'s `/import` endpoint. You'll need to adapt the following command, before either pasting it into your terminal or running it as a bash script:

```sh
sqlite3 -json -batch /path/to/your/rarbg_db.sqlite "$(cat rarbg-import.sql)" \
  | jq -r --indent 0 '.[] | . * { source: "rarbg" } | . + if .imdb != null then { contentSource: "imdb", contentId: .imdb } else {} end | del(.imdb) | del(..|nulls)' \
  | curl --verbose -H "Content-Type: application/json" -H "Connection: close" --data-binary @- http://localhost:3333/import
```

So what's happening here?

- First we are executing the SQL query we made above against the backup database; we tell SQLite to output the result as JSON. To test this bit in isolation you might try running just `sqlite3 -json -batch /path/to/your/rarbg_db.sqlite "$(cat rarbg-import.sql)"` (while testing you'll probably want to `limit` your results to say 10 or 100)
- Next we need to make some tweaks to the JSON structure, so we'll pipe the result into [jq](https://jqlang.github.io/jq/){:target="\_blank"}. You can add the line beginning `| jq` to the previous part to test what we have so far. Here we will:
  - Add a `source` field with value `rarbg`: each torrent stored in **bitmagnet** is associated with one or more sources, this association allows filtering by source within the search facility, and can carry some source-specific information such as an import ID, and numbers of seeders and leechers (more docs needed here!)
  - Add the `contentSource` and `contentId` fields which **bitmagnet** expects, containing the IMDB ID, if it exists; these are not a required field, but if you know the external IMDB or TMDB ID of your content then it will give the classifier an easier job
  - Delete the `imdb` field which won't be recognised by **bitmagnet**
  - Delete any `null` values to reduce the payload size
- Next we'll pipe the final result to **bitmagnet**'s `/import` endpoint; you'll see feedback as the import progresses; watch out for any errors in the logs!

Total time for the import will depend on the number of imported records and on your hardware. For me it took about 10 minutes to import 1.5 million records on M2 MacBook Air.

Once the import starts you should immediately start seeing the items appear in the web UI. This isn't the end of the story though; each imported item will also be sent to the classification queue to further enrich its metadata. As the queue progresses you'll start seeing more details appear in the web UI. If you're importing a large number of items, the queue can take hours to work down. Once the metadata for any given movie or TV show has been saved, we shouldn't need to query TMDB again for it, therefore the queue should accelerate as you accumulate local metadata for all the most popular content.
