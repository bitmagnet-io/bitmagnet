import { CollectionViewer, DataSource } from "@angular/cdk/collections";
import { BehaviorSubject, catchError, EMPTY, map, Observable } from "rxjs";
import * as generated from "../../graphql/generated";
import { TorrentContentService } from "./torrent-content.service";
import { GenreAgg, VideoResolutionAgg, VideoSourceAgg } from "./facet";

export class TorrentContentDataSource
  implements DataSource<generated.TorrentContent>
{
  private resultSubject = new BehaviorSubject<generated.TorrentContentResult>({
    items: [],
    totalCount: 0,
    aggregations: {},
  });
  private loadingSubject = new BehaviorSubject(false);
  private errorSubject = new BehaviorSubject<Error | undefined>(undefined);
  private lastRequestTimeSubject = new BehaviorSubject<number>(0);

  public result = this.resultSubject.asObservable();
  public items = this.result.pipe(map((result) => result.items));
  public totalCount = this.result.pipe(map((result) => result.totalCount));
  public aggregations = this.result.pipe(map((result) => result.aggregations));
  public loading = this.loadingSubject.asObservable();
  public error = this.errorSubject.asObservable();

  public torrentSourceAggs: Observable<generated.TorrentSourceAgg[]> =
    this.aggregations.pipe(map((aggs) => aggs.torrentSource ?? []));
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

  constructor(private torrentContentService: TorrentContentService) {}

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

  loadResult(input: generated.TorrentContentSearchQueryVariables): void {
    const requestTime = Date.now();
    this.loadingSubject.next(true);
    this.torrentContentService
      .search(input)
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
        }
      });
  }
}
