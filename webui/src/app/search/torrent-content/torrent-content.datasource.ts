import { CollectionViewer, DataSource } from "@angular/cdk/collections";
import { BehaviorSubject, catchError, EMPTY, map, Observable } from "rxjs";
import * as generated from "../../graphql/generated";
import { GraphQLService } from "../../graphql/graphql.service";
import { AppErrorsService } from "../../app-errors.service";
import { GenreAgg, VideoResolutionAgg, VideoSourceAgg } from "./facet";

export class TorrentContentDataSource
  implements DataSource<generated.TorrentContent>
{
  private resultSubject =
    new BehaviorSubject<generated.TorrentContentSearchResult>({
      items: [],
      totalCount: 0,
      aggregations: {},
    });
  private inputSubject =
    new BehaviorSubject<generated.TorrentContentSearchQueryVariables>({});
  private loadingSubject = new BehaviorSubject(false);
  private lastRequestTimeSubject = new BehaviorSubject<number>(0);

  public result = this.resultSubject.asObservable();
  public items = this.result.pipe(map((result) => result.items));
  public aggregations = this.result.pipe(map((result) => result.aggregations));
  public loading = this.loadingSubject.asObservable();

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

  constructor(
    private graphQLService: GraphQLService,
    private errorsService: AppErrorsService,
  ) {}

  connect({}: CollectionViewer): Observable<generated.TorrentContent[]> {
    return this.items;
  }

  disconnect(): void {
    this.resultSubject.complete();
    this.loadingSubject.complete();
  }

  loadResult(input: generated.TorrentContentSearchQueryVariables): void {
    const requestTime = Date.now();
    this.loadingSubject.next(true);
    this.inputSubject.next(input);
    this.graphQLService
      .torrentContentSearch(input)
      .pipe(
        catchError((err: Error) => {
          this.errorsService.addError(`Error loading results: ${err.message}`);
          this.loadingSubject.next(false);
          return EMPTY;
        }),
      )
      .subscribe((result) => {
        const lastRequestTime = this.lastRequestTimeSubject.getValue();
        if (requestTime > lastRequestTime) {
          this.lastRequestTimeSubject.next(requestTime);
          this.resultSubject.next(result);
          this.loadingSubject.next(false);
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
}
