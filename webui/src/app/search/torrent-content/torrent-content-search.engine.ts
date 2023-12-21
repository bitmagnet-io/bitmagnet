import { CollectionViewer, DataSource } from "@angular/cdk/collections";
import { BehaviorSubject, catchError, EMPTY, map, Observable, zip } from "rxjs";
import * as generated from "../../graphql/generated";
import { GraphQLService } from "../../graphql/graphql.service";
import { AppErrorsService } from "../../app-errors.service";
import { PageEvent } from "../../paginator/paginator.types";
import { Facet, VideoResolutionAgg, VideoSourceAgg } from "./facet";

const emptyResult: generated.TorrentContentSearchResult = {
  items: [],
  totalCount: 0,
  aggregations: {},
};

export class TorrentContentSearchEngine
  implements DataSource<generated.TorrentContent>
{
  private queryStringSubject = new BehaviorSubject<string>("");

  private contentTypeSubject = new BehaviorSubject<
    generated.ContentType | "null" | null | undefined
  >(undefined);

  private pageIndexSubject = new BehaviorSubject<number>(0);
  public pageIndex$ = this.pageIndexSubject.asObservable();

  private pageSizeSubject = new BehaviorSubject<number>(10);
  public pageSize$ = this.pageSizeSubject.asObservable();

  private pageLengthSubject = new BehaviorSubject<number>(0);
  public pageLength$ = this.pageLengthSubject.asObservable();

  private itemsResultSubject =
    new BehaviorSubject<generated.TorrentContentSearchResult>(emptyResult);
  private aggsResultSubject =
    new BehaviorSubject<generated.TorrentContentSearchResult>(emptyResult);
  private loadingSubject = new BehaviorSubject(false);
  private lastRequestTimeSubject = new BehaviorSubject<number>(0);

  public items$ = this.itemsResultSubject.pipe(map((result) => result.items));
  public aggregations$ = this.aggsResultSubject.pipe(
    map((result) => result.aggregations),
  );
  public loading$ = this.loadingSubject.asObservable();

  private torrentSourceFacet = new Facet<string, false>(
    "Torrent Source",
    "mediation",
    null,
    this.aggregations$.pipe(map((aggs) => aggs.torrentSource ?? [])),
  );
  private torrentTagFacet = new Facet<string, false>(
    "Torrent Tag",
    "sell",
    null,
    this.aggregations$.pipe(map((aggs) => aggs.torrentTag ?? [])),
  );
  private torrentFileTypeFacet = new Facet<generated.FileType, false>(
    "File Type",
    "file_present",
    null,
    this.aggregations$.pipe(map((aggs) => aggs.torrentFileType ?? [])),
  );
  private languageFacet = new Facet<generated.Language, false>(
    "Language",
    "translate",
    null,
    this.aggregations$.pipe(map((aggs) => aggs.language ?? [])),
  );
  private genreFacet = new Facet<string, false>(
    "Genre",
    "theater_comedy",
    ["movie", "tv_show"],
    this.aggregations$.pipe(map((aggs) => aggs.genre ?? [])),
  );
  private videoResolutionFacet = new Facet<generated.VideoResolution>(
    "Video Resolution",
    "screenshot_monitor",
    ["movie", "tv_show", "xxx"],
    this.aggregations$.pipe(
      map((aggs) => (aggs.videoResolution ?? []) as VideoResolutionAgg[]),
    ),
  );
  private videoSourceFacet = new Facet<generated.VideoSource>(
    "Video Source",
    "album",
    ["movie", "tv_show", "xxx"],
    this.aggregations$.pipe(
      map((aggs) => (aggs.videoSource ?? []) as VideoSourceAgg[]),
    ),
  );

  public facets: Facet<unknown, boolean>[] = [
    this.torrentSourceFacet,
    this.torrentTagFacet,
    this.torrentFileTypeFacet,
    this.languageFacet,
    this.videoResolutionFacet,
    this.videoSourceFacet,
    this.genreFacet,
  ];

  private overallTotalCountSubject = new BehaviorSubject<number>(0);
  public overallTotalCount$ = this.overallTotalCountSubject.asObservable();

  private maxTotalCountSubject = new BehaviorSubject<number>(0);
  public maxTotalCount$ = this.maxTotalCountSubject.asObservable();

  public contentTypes = contentTypes;

  constructor(
    private graphQLService: GraphQLService,
    private errorsService: AppErrorsService,
  ) {
    this.contentTypeSubject.subscribe((contentType) => {
      this.facets.forEach(
        (f) => f.isRelevant(contentType) || f.deactivateAndReset(),
      );
      this.pageIndexSubject.next(0);
      this.loadResult();
    });
  }

  connect({}: CollectionViewer): Observable<generated.TorrentContent[]> {
    return this.items$;
  }

  disconnect(): void {
    this.itemsResultSubject.complete();
    this.loadingSubject.complete();
  }

  loadResult(cached = true): void {
    const requestTime = Date.now();
    this.loadingSubject.next(true);
    const pageSize = this.pageSizeSubject.getValue();
    const queryString = this.queryStringSubject.getValue() || undefined;
    const offset = this.pageIndexSubject.getValue() * pageSize;
    const contentType = this.contentTypeSubject.getValue();
    const items = this.graphQLService
      .torrentContentSearch({
        query: {
          queryString,
          limit: pageSize,
          offset,
          hasNextPage: true,
          cached,
        },
        facets: this.facetsInput(false),
      })
      .pipe(
        catchError((err: Error) => {
          this.errorsService.addError(
            `Error loading item results: ${err.message}`,
          );
          return EMPTY;
        }),
      );
    const aggs = this.graphQLService
      .torrentContentSearch({
        query: {
          queryString,
          limit: 0,
          cached: true,
        },
        facets: this.facetsInput(true),
      })
      .pipe(
        catchError((err: Error) => {
          this.errorsService.addError(
            `Error loading aggregation results: ${err.message}`,
          );
          return EMPTY;
        }),
      );
    items.subscribe((result) => {
      const lastRequestTime = this.lastRequestTimeSubject.getValue();
      if (requestTime >= lastRequestTime) {
        this.itemsResultSubject.next(result);
      }
    });
    aggs.subscribe((result) => {
      const lastRequestTime = this.lastRequestTimeSubject.getValue();
      if (requestTime >= lastRequestTime) {
        this.aggsResultSubject.next(result);
      }
    });
    zip(items, aggs).subscribe(([i, a]) => {
      this.loadingSubject.next(false);
      this.lastRequestTimeSubject.next(requestTime);
      const overallTotalCount =
        a.aggregations.contentType
          ?.map((c) => c.count)
          .reduce((a, b) => a + b, 0) ?? 0;
      let maxTotalCount = 0;
      if (!i.hasNextPage) {
        maxTotalCount = offset + i.items.length;
      } else if (contentType === undefined) {
        maxTotalCount =
          a.aggregations.contentType
            ?.map((c) => c.count)
            .reduce((a, b) => a + b, 0) ?? 0;
      } else {
        maxTotalCount =
          a.aggregations.contentType?.find(
            (a) => (a.value ?? "null") === (contentType ?? undefined),
          )?.count ?? overallTotalCount;
      }
      this.pageLengthSubject.next(i.items.length);
      this.maxTotalCountSubject.next(maxTotalCount);
      this.overallTotalCountSubject.next(overallTotalCount);
    });
  }

  private facetsInput(aggregate: boolean): generated.TorrentContentFacetsInput {
    const contentType = this.contentTypeSubject.getValue();
    return {
      contentType: {
        aggregate,
        filter:
          contentType === "null"
            ? [null]
            : contentType
            ? [contentType]
            : undefined,
      },
      torrentSource: this.torrentSourceFacet.isActive()
        ? {
            aggregate,
            filter: aggregate
              ? undefined
              : this.torrentSourceFacet.filterValues(),
          }
        : undefined,
      torrentTag: this.torrentTagFacet.isActive()
        ? {
            aggregate,
            filter: aggregate ? undefined : this.torrentTagFacet.filterValues(),
          }
        : undefined,
      torrentFileType: this.torrentFileTypeFacet.isActive()
        ? {
            aggregate,
            filter: aggregate
              ? undefined
              : this.torrentFileTypeFacet.filterValues(),
          }
        : undefined,
      language: this.languageFacet.isActive()
        ? {
            aggregate,
            filter: aggregate ? undefined : this.languageFacet.filterValues(),
          }
        : undefined,
      genre: this.genreFacet.isActive()
        ? {
            aggregate,
            filter: aggregate ? undefined : this.genreFacet.filterValues(),
          }
        : undefined,
      videoResolution: this.videoResolutionFacet.isActive()
        ? {
            aggregate,
            filter: aggregate
              ? undefined
              : this.videoResolutionFacet.filterValues(),
          }
        : undefined,
      videoSource: this.videoSourceFacet.isActive()
        ? {
            aggregate,
            filter: aggregate
              ? undefined
              : this.videoSourceFacet.filterValues(),
          }
        : undefined,
    };
  }

  get isDeepFiltered(): boolean {
    return (
      !!this.queryStringSubject.getValue() ||
      this.facets.some((f) => f.isActive() && !f.isEmpty())
    );
  }

  selectContentType(
    contentType: generated.ContentType | "null" | null | undefined,
  ) {
    this.contentTypeSubject.next(contentType);
  }

  setQueryString(queryString: string) {
    this.queryStringSubject.next(queryString);
  }

  get hasQueryString(): boolean {
    return !!this.queryStringSubject.getValue();
  }

  firstPage() {
    this.pageIndexSubject.next(0);
  }

  handlePageEvent(event: PageEvent) {
    this.pageIndexSubject.next(event.pageIndex);
    this.pageSizeSubject.next(event.pageSize);
    this.loadResult();
  }

  contentTypeCount(type: string): Observable<number> {
    return this.aggregations$.pipe(
      map(
        (aggs) => aggs.contentType?.find((a) => a.value === type)?.count ?? 0,
      ),
    );
  }

  contentTypeInfo(type: unknown): ContentTypeInfo | undefined {
    return contentTypes[type as keyof typeof contentTypes];
  }
}

type ContentTypeInfo = {
  singular: string;
  plural: string;
  icon: string;
};

const contentTypes: Record<generated.ContentType | "null", ContentTypeInfo> = {
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
    icon: "18_up_rating_outline",
  },
  null: {
    singular: "Unknown",
    plural: "Unknown",
    icon: "question_mark",
  },
};
