---
title: Servarr Integration
parent: Tutorials
layout: default
nav_order: 1
---

# Servarr Integration

**bitmagnet**'s HTTP server exposes an endpoint at `/torznab`, allowing it to integrate with any application that supports [the Torznab specification](https://torznab.github.io/spec-1.3-draft/index.html){:target="\_blank"}, most notably apps in [the Servarr stack](https://wiki.servarr.com/){:target="\_blank"} (Prowlarr, Sonarr, Radarr etc.).

## Adding **bitmagnet** as an indexer in Prowlarr

To get started, open your Prowlarr instance, click "Add Indexer", and select "Generic Torznab" from the list.

![Prowlarr Add Indexer](/assets/images/prowlarr-1.png)

The required settings are fairly basic. Assuming you've adapted from the [example docker-compose file]({% link setup/installation.md %}#docker), and Prowlarr is on the same Docker network as **bitmagnet**, then Prowlarr should be able to access the Torznab endpoint of your **bitmagnet** instance at `http://bitmagnet:3333/torznab`. No further configuration should be needed, just click the "Test" button to ensure everything is working.

![Prowlarr configure bitmagnet](/assets/images/prowlarr-2.png)

[Depending on your Prowlarr configuration](https://wiki.servarr.com/prowlarr/settings#applications){:target="\_blank"}, the **bitmagnet** indexer should now be synced to your other \*arr applications. Alternatively, you can add **bitmagnet** as an indexer directly in those applications, following the same steps as above.
