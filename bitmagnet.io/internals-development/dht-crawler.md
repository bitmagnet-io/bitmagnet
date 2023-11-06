---
title: DHT Crawler
parent: Internals & Development
layout: default
nav_order: 2
---

# Architecture & Lifecycle of the DHT Crawler

The DHT and BitTorrent protocols are (rather impenetrably) documented at [bittorrent.org](http://bittorrent.org/beps/bep_0000.html){:target="\_blank"}. Relevant resources include:

- [BEP 5: DHT Protocol](http://bittorrent.org/beps/bep_0005.html){:target="\_blank"}
- [BEP 51: Infohash Indexing](https://www.bittorrent.org/beps/bep_0051.html){:target="\_blank"}
- [BEP 33: DHT Scrapes](https://www.bittorrent.org/beps/bep_0033.html){:target="\_blank"}
- [BEP 10: Extension Protocol](https://www.bittorrent.org/beps/bep_0010.html){:target="\_blank"}
- [The Kademlia paper](https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf){:target="\_blank"}

The rest of what I've figured out about how to implement a DHT crawler was cobbled together from [the now archived **magnetico** project](https://github.com/boramalper/magnetico){:target="\_blank"} and [anacrolix's BitTorrent libraries](https://github.com/anacrolix){:target="\_blank"}.

The following diagram illustrates roughly how the crawler has been implemented within **bitmagnet**. It's debatable if this will help stop anyone's brain from melting, including my own.

{: .warning-title }

> Todo
>
> This diagram is out-of-date and needs updating to reflect the new DHT crawler design.

```mermaid
%%{init: {"flowchart": {"defaultRenderer": "elk"}} }%%
flowchart TB
    START{Start}
    START -->STEP_OPEN_DHT_CONNECTION
    STEP_OPEN_DHT_CONNECTION(Open DHT connection)
    STEP_OPEN_DHT_CONNECTION -.->DHT
    STEP_OPEN_DHT_CONNECTION -->STEP_crawl(Crawl bootstrap nodes)
    STEP_crawl(Crawl bootstrap nodes) --> DHT_find_node[[DHT: find_node]]
    DHT_find_node -.->|Add to routing table| ROUTING_TABLE
    DHT_find_node -.->|Loop| STEP_crawl
    ROUTING_TABLE[/Routing Table/]
    STEP_select_node(Select a node from routing table and acquire lock)
    ROUTING_TABLE -.->STEP_select_node
    STEP_select_node -->DHT_sample_infohashes[[DHT: sample_infohashes]]
    DHT_sample_infohashes -->STEP_add_to_staging(Add hashes to staging)
    STEP_OPEN_DHT_CONNECTION --> STEP_select_node
    STEP_add_to_staging -->STEP_check_in_progress
    STEP_add_to_staging -->|Loop| STEP_select_node
    subgraph InfoHash staging
        STEP_check_in_progress(Is request for InfoHash already in progress?)
        STEP_check_in_progress -->|No| STEP_gather_infohashes
        STEP_gather_infohashes(Gather InfoHashes for batch DB check)
        STEP_gather_infohashes -->STEP_check_persisted_infohashes
        STEP_check_persisted_infohashes(Is InfoHash already persisted?)
        STEP_torrent_received(Torrent info received in staging)
        STEP_torrent_received -->STEP_persist_torrent
        STEP_torrent_received -->STEP_publish_classify_job
        STEP_persist_torrent(Persist torrent to database)
        STEP_publish_classify_job(Publish classify job)
        STEP_remove_torrent_from_staging("Remove torrent from staging")
        STEP_persist_torrent -->STEP_remove_torrent_from_staging
        STEP_publish_classify_job -->STEP_remove_torrent_from_staging
    end
    STEP_torrent_to_staging(Send torrent to staging)
    STEP_torrent_to_staging -->STEP_torrent_received
    STEP_remove_torrent_from_staging -->END
    STEP_persist_torrent -.->POSTGRES
    STEP_check_in_progress -->|Yes| END
    STEP_check_persisted_infohashes -->|Yes| END
    POSTGRES -.->STEP_check_persisted_infohashes
    STEP_check_persisted_infohashes -->|No| STEP_request_torrent_info(Request torrent info)
    STEP_request_torrent_info -->DHT_get_peers[[DHT: get_peers]]
    DHT_get_peers -->BT_request_meta_info[[BT: Request MetaInfo]]
    DHT_get_peers -.->|Add to routing table| ROUTING_TABLE
    STEP_request_torrent_info -->DHT_get_peers_scrape[["DHT: get_peers (scrape)"]]
    DHT_get_peers_scrape -->BT_request_meta_info[[BT: Request MetaInfo]]
    BT_request_meta_info -->STEP_meta_info_success(Did meta info request succeed for any peer?)
    STEP_meta_info_success -->|No| STEP_remove_torrent_from_staging
    STEP_meta_info_success -->|Yes| STEP_torrent_to_staging
    POSTGRES[(Postgres Database)]
    MESSAGE_QUEUE[(Redis Message Queue)]
    STEP_publish_classify_job -.->MESSAGE_QUEUE
    MESSAGE_QUEUE -.->STEP_classify_torrent(Classify torrent content)
    STEP_persist_torrent_content(Persist content metadata)
    STEP_classify_torrent -->STEP_persist_torrent_content
    STEP_persist_torrent_content -.->POSTGRES
    STEP_persist_torrent_content -->END
    DHT((DHT connection))
    DHT_find_node <-.->DHT
    DHT_sample_infohashes <-.->DHT
    DHT_get_peers <-.->DHT
    DHT_get_peers_scrape <-.->DHT
    END{End}
```

[comment]: <> (Need to enable panning and zooming for this ridiculous diagram, let the hacking commence!)
[comment]: <> (panzoom comes from https://github.com/timmywil/panzoom)

<script src="https://unpkg.com/@panzoom/panzoom@4.5.1/dist/panzoom.min.js"></script>
<script>
const i = setInterval(() => {
    const elem = document.querySelector('.language-mermaid');
    const svgElem = elem.childNodes[0];
    if (svgElem?.tagName === 'svg') {
        clearInterval(i);
        const parentElem = elem.parentElement;
        parentElem.style.overflow = 'hidden';
        const pz = Panzoom(svgElem);
        parentElem.addEventListener('wheel', pz.zoomWithWheel);
        const grandparentElem = parentElem.parentElement;
        const zoomIn = document.createElement('button');
        zoomIn.innerText = '+';
        grandparentElem.insertBefore(zoomIn, parentElem);
        zoomIn.addEventListener('click', () => pz.zoomIn());
        const zoomOut = document.createElement('button');
        zoomOut.innerText = '-';
        grandparentElem.insertBefore(zoomOut, parentElem);
        zoomOut.addEventListener('click', () => pz.zoomOut());
        const reset = document.createElement('button');
        reset.innerText = 'Reset';
        grandparentElem.insertBefore(reset, parentElem);
        reset.addEventListener('click', () => pz.reset());
    }
}, 100)
</script>
