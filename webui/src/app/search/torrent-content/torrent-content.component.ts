import { AfterContentInit, AfterViewInit, Component } from "@angular/core";
import { PageEvent } from "@angular/material/paginator";
import {
  BehaviorSubject,
  catchError,
  EMPTY,
  interval,
  startWith,
  Subscription,
  tap,
} from "rxjs";
import { FormControl } from "@angular/forms";
import {
  animate,
  state,
  style,
  transition,
  trigger,
} from "@angular/animations";
import { COMMA, ENTER } from "@angular/cdk/keycodes";
import * as generated from "../../graphql/generated";
import { GraphQLService } from "../../graphql/graphql.service";
import { AppErrorsService } from "../../app-errors.service";
import { TorrentContentDataSource } from "./torrent-content.datasource";
import { Facet } from "./facet";

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
  dataSource: TorrentContentDataSource = new TorrentContentDataSource(
    this.graphQLService,
    this.errorsService,
  );
  displayedColumns = ["summary", "size", "peers", "magnet"];
  queryString = new FormControl("");

  pageIndex = 0;
  pageSize = 10;

  result: generated.TorrentContentSearchResult = {
    items: [],
    aggregations: {},
    totalCount: 0,
  };
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

  torrentSourceFacet: Facet<string, false>;
  torrentTagFacet: Facet<string, false>;
  torrentFileTypeFacet: Facet<generated.FileType, false>;
  languageFacet: Facet<generated.Language, false>;
  genreFacet: Facet<string, false>;
  videoResolutionFacet: Facet<generated.VideoResolution>;
  videoSourceFacet: Facet<generated.VideoSource>;

  facets: Facet<unknown, boolean>[];

  autoRefresh = new FormControl<number>(0);
  autoRefreshIntervals = [0, 10, 30];
  autoRefreshSubscription: Subscription | undefined;

  readonly separatorKeysCodes = [ENTER, COMMA] as const;

  constructor(
    private graphQLService: GraphQLService,
    private errorsService: AppErrorsService,
  ) {
    this.dataSource.result.subscribe((result) => {
      this.result = result;
    });
    this.dataSource.loading.subscribe((loading) => (this.loading = loading));
    this.torrentSourceFacet = new Facet<string, false>(
      "Torrent Source",
      "mediation",
      null,
      this.dataSource.torrentSourceAggs,
    );
    this.torrentTagFacet = new Facet<string, false>(
      "Torrent Tag",
      "sell",
      null,
      this.dataSource.torrentTagAggs,
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
      null,
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
      ["movie", "tv_show", "xxx"],
      this.dataSource.videoResolutionAggs,
    );
    this.videoSourceFacet = new Facet<generated.VideoSource>(
      "Video Source",
      "album",
      ["movie", "tv_show", "xxx"],
      this.dataSource.videoSourceAggs,
    );
    this.facets = [
      this.torrentSourceFacet,
      this.torrentTagFacet,
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
        torrentTag: this.torrentTagFacet.facetParams(),
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
    const ct = this.result.aggregations.contentType;
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

  expandedItem = new (class {
    private itemSubject = new BehaviorSubject<
      generated.TorrentContent | undefined
    >(undefined);
    newTagCtrl = new FormControl<string>("");
    private editedTagsSubject = new BehaviorSubject<string[] | undefined>(
      undefined,
    );
    public readonly suggestedTags = Array<string>();
    public selectedTabIndex = 0;

    constructor(private ds: TorrentContentComponent) {
      ds.dataSource.result.subscribe((result) => {
        const item = this.itemSubject.getValue();
        if (!item) {
          return;
        }
        const nextItem = result.items.find((i) => i.id === item.id);
        this.itemSubject.next(nextItem);
      });
      this.newTagCtrl.valueChanges.subscribe((value) => {
        return ds.graphQLService
          .torrentSuggestTags({
            query: {
              prefix: value,
              exclusions: this.itemSubject.getValue()?.torrent.tagNames,
            },
          })
          .pipe(
            tap((result) => {
              this.suggestedTags.splice(
                0,
                this.suggestedTags.length,
                ...result.suggestions.map((t) => t.name),
              );
            }),
          )
          .subscribe();
      });
    }

    get id(): string | undefined {
      return this.itemSubject.getValue()?.id;
    }

    toggle(id?: string): void {
      if (id === this.id) {
        id = undefined;
      }
      const nextItem = this.ds.result.items.find((i) => i.id === id);
      const current = this.itemSubject.getValue();
      if (current?.id !== id) {
        this.itemSubject.next(nextItem);
        this.editedTagsSubject.next(undefined);
        this.newTagCtrl.reset();
        this.selectedTabIndex = 0;
      }
    }

    selectTab(index: number): void {
      this.selectedTabIndex = index;
    }

    addTag(tagName: string) {
      this.editTags((tags) => [...tags, tagName]);
      this.saveTags();
    }

    renameTag(oldTagName: string, newTagName: string) {
      this.editTags((tags) =>
        tags.map((t) => (t === oldTagName ? newTagName : t)),
      );
      this.saveTags();
    }

    deleteTag(tagName: string) {
      this.editTags((tags) => tags.filter((t) => t !== tagName));
      this.saveTags();
    }

    private editTags(fn: (tagNames: string[]) => string[]) {
      const item = this.itemSubject.getValue();
      if (!item) {
        return;
      }
      const editedTags =
        this.editedTagsSubject.getValue() ?? item.torrent.tagNames;
      this.editedTagsSubject.next(fn(editedTags));
      this.newTagCtrl.reset();
    }

    saveTags(): void {
      const expanded = this.itemSubject.getValue();
      if (!expanded) {
        return;
      }
      const editedTags = this.editedTagsSubject.getValue();
      if (!editedTags) {
        return;
      }
      this.ds.graphQLService
        .torrentSetTags({
          infoHashes: [expanded.infoHash],
          tagNames: editedTags ?? [],
        })
        .pipe(
          catchError((err: Error) => {
            this.ds.errorsService.addError(`Error saving tags: ${err.message}`);
            this.editedTagsSubject.next(undefined);
            return EMPTY;
          }),
        )
        .pipe(
          tap(() => {
            this.editedTagsSubject.next(undefined);
            this.ds.dataSource.refreshResult();
          }),
        )
        .subscribe();
    }
  })(this);
}

type ContentTypeInfo = {
  singular: string;
  plural: string;
  icon: string;
};
