---
title: Reprocess & Reclassify
parent: Tutorials
layout: default
nav_order: 4
---

# Reprocess & Reclassify Torrents

The classifier is being updated regularly, and to reclassify already-crawled torrents you'll need to run the CLI and queue them for reprocessing.

{% include callout_cli.md %}

For context: after torrents are crawled or imported, they won't show up in the UI straight away. They must first be "processed" by the job queue. This involves a few steps:

- The classifier attempts to classify the torrent (determine its content type, and match it to a known piece of content)
- The search index for the torrent is built
- (In future there's likely to be other steps here, such as running rule-based actions)
- The torrent content record is saved to the database

The `reprocess` command will re-queue torrents to allow the latest updates to be applied to their content records.

To reprocess all torrents in your index, simply run `bitmagnet reprocess`. If you've indexed a lot of torrents, this will take a while, so there are a few options available to control exactly what gets reprocessed:

- `contentType`: Only reprocess torrents of a certain content type. For example, `bitmagnet reprocess --contentType movie` will only reprocess movies. Multiple content types can be comma separated, and `null` refers to torrents of unknown content type.
- `classifyMode`: This controls how already matched torrents are handled. A torrent is "matched" if it's associated with a specific piece of content from one of the API integrations (currently only TMDB). Making a lot of API calls can take a long time, so if items are already matched you might want to just do the other processing steps without re-matching them. The available modes are:
  - `default`: Only attempt to match previously unmatched torrents
  - `rematch`: Ignore any pre-existing match and always classify from scratch

\*hints tell the classifier to use the hinted information instead of any classification results, which can save a lot of work for the classifier and help fix errors. Currently, the only way to add hints is by using [the `/import` endpoint](/tutorials/import.html).
