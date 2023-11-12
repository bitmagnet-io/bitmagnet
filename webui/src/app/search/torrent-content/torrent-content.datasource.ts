import { CollectionViewer, DataSource } from "@angular/cdk/collections";
import { BehaviorSubject, catchError, EMPTY, map, Observable, tap } from "rxjs";
import * as generated from "../../graphql/generated";
import { GraphQLService } from "../../graphql/graphql.service";
import { GenreAgg, VideoResolutionAgg, VideoSourceAgg } from "./facet";

export class TorrentContentDataSource
  implements DataSource<generated.TorrentContent>
{
  private resultSubject = new BehaviorSubject<generated.TorrentContentResult>({
    items: [],
    totalCount: 0,
    aggregations: {},
  });
  private inputSubject =
    new BehaviorSubject<generated.SearchTorrentContentQueryVariables>({});
  private loadingSubject = new BehaviorSubject(false);
  private errorSubject = new BehaviorSubject<Error | undefined>(undefined);
  private lastRequestTimeSubject = new BehaviorSubject<number>(0);
  private expandedItemSubject = new BehaviorSubject<
    generated.TorrentContent | undefined
  >(undefined);
  private editedTagsSubject = new BehaviorSubject<string[] | undefined>(
    undefined,
  );

  public result = this.resultSubject.asObservable();
  public items = this.result.pipe(map((result) => result.items));
  public totalCount = this.result.pipe(map((result) => result.totalCount));
  public aggregations = this.result.pipe(map((result) => result.aggregations));
  public loading = this.loadingSubject.asObservable();
  public error = this.errorSubject.asObservable();

  public torrentSourceAggs: Observable<generated.TorrentSourceAgg[]> =
    this.aggregations.pipe(map((aggs) => aggs.torrentSource ?? []));
  public torrentTagAggs: Observable<generated.TorrentTagAgg[]> =
    this.aggregations.pipe(map((aggs) => aggs.torrentTag ?? []));
  public torrentFileTypeAggs: Observable<generated.TorrentFileTypeAgg[]> =
    this.aggregations.pipe(map((aggs) => aggs.torrentFileType ?? []));
  public languageAggs: Observable<generated.LanguageAgg[]> =
    this.aggregations.pipe(map((aggs) => aggs.language ?? []));
  public genreAggs: Observable<GenreAgg[]> = this.aggregations.pipe(
    map((aggs) => aggs.genre ?? []),
  );
  public videoResolutionAggs: Observable<VideoResolutionAgg[]> =
    this.aggregations.pipe(
      map((aggs) => (aggs.videoResolution ?? []) as VideoResolutionAgg[]),
    );
  public videoSourceAggs: Observable<VideoSourceAgg[]> = this.aggregations.pipe(
    map((aggs) => (aggs.videoSource ?? []) as VideoSourceAgg[]),
  );

  constructor(private graphQLService: GraphQLService) {}

  connect(
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    _: CollectionViewer,
  ): Observable<generated.TorrentContent[]> {
    return this.items;
  }

  disconnect(): void {
    this.resultSubject.complete();
    this.loadingSubject.complete();
    this.errorSubject.complete();
  }

  loadResult(input: generated.SearchTorrentContentQueryVariables): void {
    const requestTime = Date.now();
    this.loadingSubject.next(true);
    this.inputSubject.next(input);
    this.graphQLService
      .searchTorrentContent(input)
      .pipe(
        catchError((err: Error) => {
          this.errorSubject.next(err);
          this.loadingSubject.next(false);
          return EMPTY;
        }),
      )
      .subscribe((result) => {
        const lastRequestTime = this.lastRequestTimeSubject.getValue();
        if (requestTime > lastRequestTime) {
          this.lastRequestTimeSubject.next(requestTime);
          this.errorSubject.next(undefined);
          this.resultSubject.next(result);
          this.loadingSubject.next(false);
          const expandedItem = this.expandedItemSubject.getValue();
          if (
            expandedItem &&
            !result.items.some((i) => i.id === expandedItem.id)
          ) {
            this.expandItem(undefined);
          }
        }
      });
  }

  refreshResult(): void {
    const input = this.inputSubject.getValue();
    const uncachedInput = {
      ...input,
      query: {
        ...input.query,
        cached: false,
      },
    };
    return this.loadResult(uncachedInput);
  }

  expandItem(id?: string): void {
    const nextItem = this.resultSubject
      .getValue()
      .items.find((i) => i.id === id);
    const current = this.expandedItemSubject.getValue();
    if (current?.id !== id) {
      this.expandedItemSubject.next(nextItem);
      this.editedTagsSubject.next(undefined);
    }
  }

  expandedSaveTags(): void {
    const expanded = this.expandedItemSubject.getValue();
    if (!expanded) {
      return;
    }
    const editedTags = this.editedTagsSubject.getValue();
    if (!editedTags) {
      return;
    }
    this.graphQLService
      .torrentSetTags({
        infoHashes: [expanded.infoHash],
        tagNames: editedTags,
      })
      .pipe(
        catchError((err: Error) => {
          this.errorSubject.next(err);
          return EMPTY;
        }),
      )
      .pipe(
        tap(() => {
          this.editedTagsSubject.next(undefined);
          this.refreshResult();
        }),
      );
  }
}
