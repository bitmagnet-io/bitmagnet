import { CollectionViewer, DataSource } from "@angular/cdk/collections";
import { BehaviorSubject, catchError, EMPTY, map, Observable } from "rxjs";
import * as generated from "../graphql/generated";
import { GraphQLService } from "../graphql/graphql.service";
import { AppErrorsService } from "../app-errors.service";
import { PageEvent } from "../paginator/paginator.types";
import { Facet, VideoResolutionAgg, VideoSourceAgg } from "./facet";

const emptyResult: generated.TorrentContentSearchResult = {
  items: [],
  totalCount: 0,
  totalCountIsEstimate: false,
  aggregations: {},
};

type BudgetedCount = {
  count: number;
  isEstimate: boolean;
};

const emptyBudgetedCount = {
  count: 0,
  isEstimate: false,
};

export class SearchEngine implements DataSource<generated.TorrentContent> {
  private queryStringSubject = new BehaviorSubject<string>("");

  private contentTypeSubject = new BehaviorSubject<
    generated.ContentType | "null" | null | undefined
  >(undefined);

  orderBySubject = new BehaviorSubject<OrderBySelection>({
    field: "PublishedAt",
    descending: true,
  });

  orderByOptions = orderByOptions;

  private pageIndexSubject = new BehaviorSubject<number>(0);
  public pageIndex$ = this.pageIndexSubject.asObservable();

  private pageSizeSubject = new BehaviorSubject<number>(20);
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
  public hasNextPage$ = this.itemsResultSubject.pipe(
    map((result) => result.hasNextPage),
  );

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

  private overallTotalCountSubject = new BehaviorSubject<BudgetedCount>(
    emptyBudgetedCount,
  );
  public overallTotalCount$ = this.overallTotalCountSubject.asObservable();

  private totalCountSubject = new BehaviorSubject<BudgetedCount>(
    emptyBudgetedCount,
  );
  public totalCount$ = this.totalCountSubject.asObservable();

  public contentTypes = contentTypes;
  public availableContentTypes = new Set<string>();

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
    const orderBy = this.orderBySubject.getValue();
    const items = this.graphQLService
      .torrentContentSearch({
        query: {
          queryString,
          limit: pageSize,
          offset,
          hasNextPage: true,
          cached,
          totalCount: true,
        },
        facets: this.facetsInput(true),
        orderBy: [orderBy],
      })
      .pipe(
        catchError((err: Error) => {
          this.errorsService.addError(
            `Error loading item results: ${err.message}`,
          );
          return EMPTY;
        }),
      );
    items.subscribe((result) => {
      const lastRequestTime = this.lastRequestTimeSubject.getValue();
      if (requestTime >= lastRequestTime) {
        this.itemsResultSubject.next(result);
        this.aggsResultSubject.next(result);
        this.loadingSubject.next(false);
        this.lastRequestTimeSubject.next(requestTime);
        this.pageLengthSubject.next(result.items.length);
        this.totalCountSubject.next({
          count: result.totalCount,
          isEstimate: result.totalCountIsEstimate,
        });
        let overallTotalCount = 0;
        let overallIsEstimate = false;
        for (const ct of result.aggregations.contentType ?? []) {
          overallTotalCount += ct.count;
          overallIsEstimate = overallIsEstimate || ct.isEstimate;
          this.availableContentTypes.add(ct.value ?? "null");
        }
        this.overallTotalCountSubject.next({
          count: overallTotalCount,
          isEstimate: overallIsEstimate,
        });
      }
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
      torrentSource: facetInput(this.torrentSourceFacet, aggregate),
      torrentTag: facetInput(this.torrentTagFacet, aggregate),
      torrentFileType: facetInput(this.torrentFileTypeFacet, aggregate),
      language: facetInput(this.languageFacet, aggregate),
      genre: facetInput(this.genreFacet, aggregate),
      videoResolution: facetInput(this.videoResolutionFacet, aggregate),
      videoSource: facetInput(this.videoSourceFacet, aggregate),
    };
  }

  selectContentType(
    contentType: generated.ContentType | "null" | null | undefined,
  ) {
    this.contentTypeSubject.next(contentType);
  }

  setQueryString(queryString: string) {
    if (this.queryStringSubject.getValue() === queryString) {
      return;
    }
    this.queryStringSubject.next(queryString);
    if (queryString) {
      this.orderBySubject.next({
        field: "Relevance",
        descending: true,
      });
    } else {
      this.orderBySubject.next({
        field: "PublishedAt",
        descending: true,
      });
    }
    this.firstPage();
    this.loadResult();
  }

  selectOrderBy(field: generated.TorrentContentOrderBy) {
    this.orderBySubject.next({
      field,
      descending: this.orderByOptions[field]?.descending ?? false,
    });
    this.loadResult();
  }

  toggleOrderByDirection() {
    const value = this.orderBySubject.getValue();
    this.orderBySubject.next({
      field: value.field,
      descending: !value.descending,
    });
    this.loadResult();
  }

  firstPage() {
    this.pageIndexSubject.next(0);
  }

  handlePageEvent(event: PageEvent) {
    this.pageIndexSubject.next(event.pageIndex);
    this.pageSizeSubject.next(event.pageSize);
    this.loadResult();
  }

  contentTypeCount(type: string): Observable<{
    count: number;
    isEstimate: boolean;
  }> {
    return this.aggregations$.pipe(
      map((aggs) => {
        const agg = aggs.contentType?.find((a) => (a.value ?? "null") === type);
        return agg ?? { count: 0, isEstimate: false };
      }),
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
  ebook: {
    singular: "E-Book",
    plural: "E-Books",
    icon: "auto_stories",
  },
  comic: {
    singular: "Comic",
    plural: "Comics",
    icon: "comic_bubble",
  },
  audiobook: {
    singular: "Audiobook",
    plural: "Audiobooks",
    icon: "mic",
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

function facetInput<T = unknown, _allowNull extends boolean = true>(
  facet: Facet<T, _allowNull>,
  aggregate: boolean,
) {
  return facet.isActive()
    ? {
        aggregate,
        filter: facet.filterValues(),
      }
    : undefined;
}

const orderByOptions: Partial<
  Record<generated.TorrentContentOrderBy, OrderByInfo>
> = {
  Relevance: {
    label: "Relevance",
    descending: true,
  },
  PublishedAt: {
    label: "Published",
    descending: true,
  },
  UpdatedAt: {
    label: "Updated",
    descending: true,
  },
  Size: {
    label: "Size",
    descending: true,
  },
  Files: {
    label: "Files Count",
    descending: true,
  },
  Seeders: {
    label: "Seeders",
    descending: true,
  },
  Leechers: {
    label: "Leechers",
    descending: true,
  },
  Name: {
    label: "Name",
    descending: false,
  },
};

type OrderByInfo = {
  label: string;
  descending: boolean;
};

type OrderBySelection = {
  field: generated.TorrentContentOrderBy;
  descending: boolean;
};
