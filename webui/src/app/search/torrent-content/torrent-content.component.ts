import { AfterContentInit, AfterViewInit, Component } from "@angular/core";
import { PageEvent } from "@angular/material/paginator";
import { interval, startWith, Subscription, tap } from "rxjs";
import { FormControl } from "@angular/forms";
import {
  animate,
  state,
  style,
  transition,
  trigger,
} from "@angular/animations";
import { MatSnackBar } from "@angular/material/snack-bar";
import * as generated from "../../graphql/generated";
import { TorrentContentService } from "./torrent-content.service";
import { TorrentContentDataSource } from "./torrent-content.datasource";
import { Facet } from "./facet";

type ContentTypeInfo = {
  singular: string;
  plural: string;
  icon: string;
};

@Component({
  selector: "app-torrent-content",
  templateUrl: "./torrent-content.component.html",
  styleUrls: ["./torrent-content.component.scss"],
  animations: [
    trigger("detailExpand", [
      state("collapsed", style({ height: "0px", minHeight: "0" })),
      state("expanded", style({ height: "*" })),
      transition(
        "expanded <=> collapsed",
        animate("225ms cubic-bezier(0.4, 0.0, 0.2, 1)"),
      ),
    ]),
  ],
})
export class TorrentContentComponent
  implements AfterContentInit, AfterViewInit
{
  dataSource: TorrentContentDataSource;
  displayedColumns = ["summary", "size", "peers", "magnet"];
  queryString = new FormControl("");

  pageIndex = 0;
  pageSize = 10;

  aggregations: generated.TorrentContentResult["aggregations"] = {};
  totalCount = 0;
  loading = true;

  contentTypes: Record<generated.ContentType | "null", ContentTypeInfo> = {
    movie: {
      singular: "Movie",
      plural: "Movies",
      icon: "movie",
    },
    tv_show: {
      singular: "TV Show",
      plural: "TV Shows",
      icon: "live_tv",
    },
    music: {
      singular: "Music",
      plural: "Music",
      icon: "music_note",
    },
    book: {
      singular: "Book",
      plural: "Books",
      icon: "auto_stories",
    },
    software: {
      singular: "Software",
      plural: "Software",
      icon: "desktop_windows",
    },
    game: {
      singular: "Game",
      plural: "Games",
      icon: "sports_esports",
    },
    xxx: {
      singular: "XXX",
      plural: "XXX",
      icon: "18_up_rating",
    },
    null: {
      singular: "Unknown",
      plural: "Unknown",
      icon: "question_mark",
    },
  };
  contentTypeEntries = Object.entries(this.contentTypes).map(
    ([type, info]) => ({
      type: type as generated.ContentType,
      ...info,
    }),
  );

  contentType = new FormControl<generated.ContentType | "null" | undefined>(
    undefined,
  );

  expandedItem: string | null = null;

  torrentSourceFacet: Facet<string, false>;
  torrentFileTypeFacet: Facet<generated.FileType, false>;
  languageFacet: Facet<generated.Language, false>;
  genreFacet: Facet<string, false>;
  videoResolutionFacet: Facet<generated.VideoResolution>;
  videoSourceFacet: Facet<generated.VideoSource>;

  facets: Facet<unknown, boolean>[];

  autoRefresh = new FormControl<number>(0);
  autoRefreshIntervals = [0, 10, 30];
  autoRefreshSubscription: Subscription | undefined;

  constructor(
    private torrentContentService: TorrentContentService,
    private errorSnackBar: MatSnackBar,
  ) {
    this.dataSource = new TorrentContentDataSource(this.torrentContentService);
    this.dataSource.error.subscribe((error: Error | undefined) => {
      if (error) {
        this.errorSnackBar.open(
          `Error fetching results: ${error.message}`,
          "Dismiss",
          {
            panelClass: ["snack-bar-error"],
          },
        );
      } else {
        this.errorSnackBar.dismiss();
      }
    });
    this.dataSource.items.subscribe((items) => {
      if (this.expandedItem && !items.some((i) => i.id === this.expandedItem)) {
        this.expandedItem = null;
      }
    });
    this.dataSource.aggregations.subscribe((aggregations) => {
      this.aggregations = aggregations;
    });
    this.dataSource.totalCount.subscribe(
      (totalCount) => (this.totalCount = totalCount),
    );
    this.dataSource.loading.subscribe((loading) => (this.loading = loading));
    this.torrentSourceFacet = new Facet<string, false>(
      "Torrent Source",
      "mediation",
      null,
      this.dataSource.torrentSourceAggs,
    );
    this.torrentFileTypeFacet = new Facet<generated.FileType, false>(
      "File Type",
      "file_present",
      null,
      this.dataSource.torrentFileTypeAggs,
    );
    this.languageFacet = new Facet<generated.Language, false>(
      "Language",
      "translate",
      ["movie", "tv_show"],
      this.dataSource.languageAggs,
    );
    this.genreFacet = new Facet<string, false>(
      "Genre",
      "theater_comedy",
      ["movie", "tv_show"],
      this.dataSource.genreAggs,
    );
    this.videoResolutionFacet = new Facet<generated.VideoResolution>(
      "Video Resolution",
      "screenshot_monitor",
      ["movie", "tv_show"],
      this.dataSource.videoResolutionAggs,
    );
    this.videoSourceFacet = new Facet<generated.VideoSource>(
      "Video Source",
      "album",
      ["movie", "tv_show"],
      this.dataSource.videoSourceAggs,
    );
    this.facets = [
      this.torrentSourceFacet,
      this.torrentFileTypeFacet,
      this.languageFacet,
      this.videoResolutionFacet,
      this.videoSourceFacet,
      this.genreFacet,
    ];
  }

  ngAfterContentInit() {
    this.loadResult();
  }

  ngAfterViewInit() {
    this.autoRefresh.valueChanges.subscribe((n) => {
      if (this.autoRefreshSubscription) {
        this.autoRefreshSubscription.unsubscribe();
      }
      if (n) {
        this.autoRefreshSubscription = interval(n * 1000)
          .pipe(
            startWith(0),
            tap(() => this.loadResult(false)),
          )
          .subscribe();
      }
    });
    this.contentType.valueChanges.subscribe((value) => {
      this.facets.forEach((f) => f.isRelevant(value) || f.deactivateAndReset());
      this.pageIndex = 0;
      this.loadResult();
    });
  }

  loadResult(cached = true) {
    this.dataSource.loadResult({
      query: {
        queryString:
          this.queryString.value === "" ? null : this.queryString.value,
        limit: this.pageSize,
        offset: this.pageIndex * this.pageSize,
        cached,
      },
      facets: {
        contentType: {
          aggregate: true,
          filter: (() => {
            const value = this.contentType.value;
            if (value == null) {
              return undefined;
            }
            return value === "null" ? [null] : [value];
          })(),
        },
        torrentSource: this.torrentSourceFacet.facetParams(),
        torrentFileType: this.torrentFileTypeFacet.facetParams(),
        language: this.languageFacet.facetParams(),
        genre: this.genreFacet.facetParams(),
        videoResolution: this.videoResolutionFacet.facetParams(),
        videoSource: this.videoSourceFacet.facetParams(),
      },
    });
  }

  handlePageEvent(event: PageEvent) {
    this.pageIndex = event.pageIndex;
    this.pageSize = event.pageSize;
    this.loadResult();
  }

  addTorrent(magnet: string, category: string | undefined ) {
    var uri = ""
    if (category == undefined) {
      uri = '/add/'+ magnet + '/other'
    } else {
      uri = '/add/'+ magnet + '/' + category.toLowerCase()
    }
    fetch(uri );
  }

  /**
   * Workaround for untyped table cell definitions
   */
  item(item: generated.TorrentContent): generated.TorrentContent {
    return item;
  }

  contentTypeInfo(type?: generated.ContentType | "null" | null) {
    return type ? this.contentTypes[type] : undefined;
  }

  contentTypeAgg(v: unknown) {
    const ct = this.aggregations.contentType;
    if (ct) {
      return ct.find((a) => (a.value ?? "null") === v);
    }
    return undefined;
  }

  getAttribute(
    item: generated.TorrentContent,
    key: string,
    source?: string,
  ): string | undefined {
    return item.content?.attributes?.find(
      (a) => a.key === key && (source === undefined || a.source === source),
    )?.value;
  }

  getCollections(
    item: generated.TorrentContent,
    type: string,
  ): string[] | undefined {
    const collections = item.content?.collections
      ?.filter((a) => a.type === type)
      .map((a) => a.name);
    return collections?.length ? collections.sort() : undefined;
  }
}
