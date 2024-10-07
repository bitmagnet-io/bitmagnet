import {
  BehaviorSubject,
  catchError,
  debounceTime,
  EMPTY,
  Observable,
} from 'rxjs';
import { Apollo } from 'apollo-angular';
import { map } from 'rxjs/operators';
import * as generated from '../graphql/generated';
import { ErrorsService } from '../errors/errors.service';
import {
  autoRefreshIntervals,
  durationSeconds,
  emptyParams,
  emptyResult,
  timeframeLengths,
} from './torrent-metrics.constants';
import {
  BucketParams,
  EventName,
  EventBucket,
  EventBucketEntries,
  EventBuckets,
  Params,
  TorrentEvents,
  Result,
  TimeframeName,
  AutoRefreshInterval,
  BucketSpan,
} from './torrent-metrics.types';

export class TorrentMetricsController {
  private paramsSubject: BehaviorSubject<Params>;
  public params$: Observable<Params>;
  private variablesSubject: BehaviorSubject<generated.TorrentMetricsQueryVariables>;
  private rawResultSubject = new BehaviorSubject<generated.TorrentMetricsQuery>(
    {
      torrent: {
        metrics: {
          buckets: [],
        },
      },
    },
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
    if (d === 'AUTO') {
      return 'hour';
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
        multiplier: multiplier ?? 'AUTO',
      },
    }));
  }

  setBucketMultiplier(multiplier: number | 'AUTO') {
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
        fetchPolicy: 'no-cache',
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
      params.buckets.duration === 'AUTO' ? 'hour' : params.buckets.duration,
    sources: params.source ? [params.source] : undefined,
    startTime: new Date(
      new Date().getTime() - 1000 * timeframeLengths[params.buckets.timeframe],
    ).toISOString(),
  },
});

const fromEntries = <K extends string, V>(
  entries: Array<[K, V]>,
): Partial<Record<K, V>> =>
  Object.fromEntries(entries) as Partial<Record<K, V>>;

const createResult = (
  params: Params,
  rawResult: generated.TorrentMetricsQuery,
): Result => {
  const { bucketParams, earliestBucket } = createBucketParams(
    params,
    rawResult,
  );
  const sources = Object.entries(
    rawResult.torrent.metrics.buckets.reduce<
      Record<string, Partial<Record<EventName, EventBucketEntries>>>
    >((acc, next) => {
      if (next.source !== (params.source ?? next.source)) {
        return acc;
      }
      let bucket: NormalizedBucket | undefined = normalizeBucket(
        next.bucket,
        bucketParams,
      );
      if (earliestBucket && earliestBucket.index > bucket.index) {
        bucket = undefined;
      }
      if (!bucket) {
        return acc;
      }
      const currentEventBuckets = acc[next.source] ?? [];
      return {
        ...acc,
        [next.source]: {
          created: !next.updated
            ? {
                ...currentEventBuckets.created,
                [bucket.key]: {
                  count:
                    next.count +
                    (currentEventBuckets.created?.[bucket.key]?.count ?? 0),
                  startTime: bucket.start,
                },
              }
            : currentEventBuckets.created,
          updated: next.updated
            ? {
                ...currentEventBuckets.updated,
                [bucket.key]: {
                  count:
                    next.count +
                    (currentEventBuckets.updated?.[bucket.key]?.count ?? 0),
                  startTime: bucket.start,
                },
              }
            : currentEventBuckets.updated,
        },
      };
    }, {}),
  ).map(([source, eventBuckets]) => {
    let events: TorrentEvents | undefined;
    // const bucketKeys = Object.keys(eventBuckets).sort()
    if (Object.keys(eventBuckets).length) {
      const bucketDates = Array<number>();
      const buckets: EventBuckets = fromEntries(
        Array<EventName>('created', 'updated').flatMap<
          [EventName, EventBucket]
        >((event): [EventName, EventBucket][] => {
          const entries = fromEntries(
            Object.entries(eventBuckets[event] ?? {})
              .filter(([, v]) => v?.count)
              .sort(([a], [b]) => (parseInt(a) < parseInt(b) ? 1 : -1)),
          );
          const keys = Object.keys(entries);
          if (!keys.length) {
            return [];
          }
          const earliestBucket = parseInt(keys[0]);
          const latestBucket = parseInt(keys[keys.length - 1]);
          bucketDates.push(earliestBucket, latestBucket);
          return [
            [
              event,
              {
                earliestBucket,
                latestBucket,
                entries,
              },
            ],
          ];
        }),
      );
      bucketDates.sort();
      events = {
        bucketDuration: bucketParams.duration,
        earliestBucket: bucketDates[0],
        latestBucket: bucketDates[bucketDates.length - 1],
        eventBuckets: buckets,
      };
    }
    return {
      source,
      events,
      isEmpty: !events?.eventBuckets,
    };
  });
  let bucketSpan: BucketSpan | undefined;
  const earliestFoundBucket = sources
    .flatMap((q) => (q.events ? [q.events.earliestBucket] : []))
    .sort()[0];
  const latestFoundBucket = sources
    .flatMap((q) => (q.events ? [q.events.latestBucket] : []))
    .sort()
    .reverse()[0];
  if (earliestFoundBucket && latestFoundBucket) {
    bucketSpan = {
      earliestBucket: earliestFoundBucket,
      latestBucket: latestFoundBucket,
    };
  }
  return {
    params: {
      ...params,
      buckets: bucketParams,
    },
    sources,
    bucketSpan,
  };
};

const createBucketParams = (
  params: Params,
  rawResult: generated.TorrentMetricsQuery,
): {
  bucketParams: BucketParams<false>;
  earliestBucket: NormalizedBucket;
  latestBucket: NormalizedBucket;
} => {
  const duration =
    params.buckets.duration === 'AUTO' ? 'hour' : params.buckets.duration;
  let multiplier =
    params.buckets.multiplier === 'AUTO' ? 1 : params.buckets.multiplier;
  const timeframe = params.buckets.timeframe;
  const now = new Date();
  const nowBucket = normalizeBucket(now, { duration, multiplier });
  const startBucket = normalizeBucket(
    now.getTime() - 1000 * timeframeLengths[timeframe],
    {
      duration,
      multiplier,
    },
  );
  const allBuckets = [
    startBucket,
    ...rawResult.torrent.metrics.buckets.flatMap((b) => [
      normalizeBucket(b.bucket, { duration, multiplier }),
    ]),
    nowBucket,
  ]
    .filter((b) => b.index >= startBucket.index)
    .sort((a, b) => a.index - b.index);
  const minBucket = allBuckets[0];
  const maxBucket = allBuckets[allBuckets.length - 1];
  if (params.buckets.multiplier === 'AUTO') {
    const targetSpan = 20;
    const span = maxBucket.index - minBucket.index;
    multiplier = Math.min(
      60,
      Math.max(Math.floor(span / (targetSpan * 5)) * 5, 1),
    );
  }
  return {
    bucketParams: {
      duration,
      multiplier,
      timeframe,
    },
    earliestBucket: normalizeBucket(
      now.getTime() - 1000 * timeframeLengths[timeframe],
      {
        duration,
        multiplier,
      },
    ),
    latestBucket: normalizeBucket(
      Math.max(now.getTime(), maxBucket.start.getTime()),
      { duration, multiplier },
    ),
  };
};

type NormalizedBucket = {
  key: string;
  index: number;
  start: Date;
};

export const normalizeBucket = (
  rawDate: string | number | Date,
  params: Pick<BucketParams<false>, 'duration' | 'multiplier'>,
): NormalizedBucket => {
  const date = new Date(rawDate);
  const msMultiplier =
    1000 * durationSeconds[params.duration] * params.multiplier;
  const baseNumber = Math.floor(date.getTime() / msMultiplier);
  return {
    key: `${baseNumber}`,
    index: baseNumber,
    start: new Date(baseNumber * msMultiplier),
  };
};
