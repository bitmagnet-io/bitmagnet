import {
  BehaviorSubject,
  catchError,
  debounceTime,
  EMPTY,
  Observable,
} from "rxjs";
import { Apollo } from "apollo-angular";
import { map } from "rxjs/operators";
import * as generated from "../../graphql/generated";
import { ErrorsService } from "../../errors/errors.service";
import {
  autoRefreshIntervals,
  emptyParams,
  emptyRawResult,
  emptyResult,
  timeframeLengths,
} from "./torrent-metrics.constants";
import {
  EventName,
  Params,
  Result,
  TimeframeName,
  AutoRefreshInterval,
} from "./torrent-metrics.types";
import { createResult } from "./torrent-metrics.utils";

export class TorrentMetricsController {
  private paramsSubject: BehaviorSubject<Params>;
  public params$: Observable<Params>;
  private variablesSubject: BehaviorSubject<generated.TorrentMetricsQueryVariables>;
  private rawResultSubject = new BehaviorSubject<generated.TorrentMetricsQuery>(
    emptyRawResult,
  );
  private resultSubject = new BehaviorSubject<Result>(emptyResult);
  public result$ = this.resultSubject.asObservable();
  private loadingSubject = new BehaviorSubject(false);

  private refreshTimeout?: number;

  constructor(
    private apollo: Apollo,
    initParams: Params = emptyParams,
    private errorsService: ErrorsService,
  ) {
    this.paramsSubject = new BehaviorSubject<Params>(initParams);
    this.params$ = this.paramsSubject.asObservable();
    this.variablesSubject =
      new BehaviorSubject<generated.TorrentMetricsQueryVariables>(
        createVariables(initParams),
      );
    this.paramsSubject.pipe(debounceTime(50)).subscribe((params) => {
      const variables = this.variablesSubject.getValue();
      const nextVariables = createVariables(params);
      if (JSON.stringify(variables) !== JSON.stringify(nextVariables)) {
        this.variablesSubject.next(nextVariables);
      } else {
        this.resultSubject.next(
          createResult(params, this.rawResultSubject.getValue()),
        );
      }
    });
    this.variablesSubject
      .pipe(debounceTime(50))
      .subscribe((variables) => this.request(variables));
    this.rawResultSubject.subscribe((rawResult) => {
      const params = this.paramsSubject.getValue();
      this.resultSubject.next(createResult(params, rawResult));
      this.setInterval(params.autoRefresh);
    });
  }

  private setInterval(interval?: AutoRefreshInterval) {
    clearTimeout(this.refreshTimeout);
    const delay = autoRefreshIntervals[interval ?? this.params.autoRefresh];
    if (delay) {
      this.refreshTimeout = setTimeout(() => {
        this.refresh();
      }, delay * 1000);
    }
  }

  get params(): Params {
    return this.paramsSubject.getValue();
  }

  get bucketDuration(): generated.MetricsBucketDuration {
    const d = this.params.buckets.duration;
    if (d === "AUTO") {
      return "hour";
    }
    return d;
  }

  get bucketMultiplier(): number {
    return (
      this.resultSubject.getValue().params.buckets.multiplier ??
      this.params.buckets.multiplier
    );
  }

  get loading(): boolean {
    return this.loadingSubject.getValue();
  }

  setTimeframe(timeframe: TimeframeName) {
    this.updateParams((p) => ({
      ...p,
      buckets: {
        ...p.buckets,
        timeframe,
      },
    }));
  }

  setSource(source: string | null) {
    this.updateParams((p) => ({
      ...p,
      source: source ?? undefined,
    }));
  }

  setBucketDuration(
    duration: generated.MetricsBucketDuration,
    multiplier?: number,
  ) {
    this.updateParams((p) => ({
      ...p,
      buckets: {
        ...p.buckets,
        duration,
        multiplier: multiplier ?? "AUTO",
      },
    }));
  }

  setBucketMultiplier(multiplier: number | "AUTO") {
    this.updateParams((p) => ({
      ...p,
      buckets: {
        ...p.buckets,
        multiplier,
      },
    }));
  }

  setEvent(event: EventName | null) {
    this.updateParams((p) => ({
      ...p,
      event: event ?? undefined,
    }));
  }

  setAutoRefreshInterval(autoRefreshInterval: AutoRefreshInterval) {
    this.updateParams((p) => ({
      ...p,
      autoRefresh: autoRefreshInterval,
    }));
  }

  private updateParams(fn: (p: Params) => Params) {
    this.paramsSubject.next(fn(this.params));
  }

  refresh() {
    this.variablesSubject.next(this.variablesSubject.getValue());
  }

  private request(variables: generated.TorrentMetricsQueryVariables) {
    clearTimeout(this.refreshTimeout);
    this.loadingSubject.next(true);
    return this.apollo
      .query<
        generated.TorrentMetricsQuery,
        generated.TorrentMetricsQueryVariables
      >({
        query: generated.TorrentMetricsDocument,
        variables,
        fetchPolicy: "no-cache",
      })
      .pipe(
        map((r) => {
          if (r) {
            this.loadingSubject.next(false);
            this.rawResultSubject.next(r.data);
          }
        }),
      )
      .pipe(
        catchError((err: Error) => {
          this.errorsService.addError(
            `Failed to load torrent metrics: ${err.message}`,
          );
          this.loadingSubject.next(false);
          this.setInterval();
          return EMPTY;
        }),
      )
      .subscribe();
  }
}

const createVariables = (
  params: Params,
): generated.TorrentMetricsQueryVariables => ({
  input: {
    bucketDuration:
      params.buckets.duration === "AUTO" ? "hour" : params.buckets.duration,
    sources: params.source ? [params.source] : undefined,
    startTime: new Date(
      new Date().getTime() - 1000 * timeframeLengths[params.buckets.timeframe],
    ).toISOString(),
  },
});
