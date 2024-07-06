import {BehaviorSubject, debounce, debounceTime} from "rxjs";
import {
  BucketParams,
  EventName, EventBucket,
  EventBucketEntries,
  EventBuckets,
  Params,
  QueueEvents,
  Result,
  StatusCounts, TimeframeName, AutoRefreshInterval
} from "./queue-metrics.types";
import {
  autoRefreshIntervals,
  durationSeconds,
  emptyParams,
  emptyResult,
  emptyStatusCounts,
  timeframeLengths
} from "./queue-metrics.constants";
import {Apollo} from "apollo-angular";
import * as generated from "../graphql/generated";
import {map} from "rxjs/operators";
import {parse as parseDuration, toSeconds} from "iso8601-duration";
import {QueueMetricsBucketDuration} from "../graphql/generated";

export class QueueMetricsController {
  private paramsSubject :BehaviorSubject<Params>;
  private variablesSubject: BehaviorSubject<generated.QueueMetricsQueryVariables>
  private rawResultSubject = new BehaviorSubject<generated.QueueMetricsQuery>({
    queue: {
      metrics: []
    }
  })
  private resultSubject = new BehaviorSubject<Result>(emptyResult)
  public result$ = this.resultSubject.asObservable()
  private loadingSubject = new BehaviorSubject(false);
  public loading$ = this.loadingSubject.asObservable();

  private refreshTimeout?: number;

  constructor(
    private apollo: Apollo,
    initParams: Params = emptyParams
  ) {
    this.paramsSubject = new BehaviorSubject<Params>(initParams);
    this.variablesSubject = new BehaviorSubject<generated.QueueMetricsQueryVariables>(createVaraibles(initParams))
    this.paramsSubject.subscribe((params) => {
      const variables = this.variablesSubject.getValue();
      const nextVariables = createVaraibles(params);
      if (JSON.stringify(variables) !== JSON.stringify(nextVariables)) {
        this.variablesSubject.next(nextVariables);
      } else {
        this.resultSubject.next(createResult(params, this.rawResultSubject.getValue()))
      }
    })
    this.variablesSubject.pipe(debounceTime(100)).subscribe((variables) =>
      this.request(variables)
    )
    this.rawResultSubject.subscribe((rawResult) => {
      const params = this.paramsSubject.getValue();
      this.resultSubject.next(createResult(params, rawResult))
      this.setInterval(params.autoRefresh)
    })
  }

  private setInterval(interval: AutoRefreshInterval) {
    clearTimeout(this.refreshTimeout);
    const delay = autoRefreshIntervals[interval];
    if (delay) {
      this.refreshTimeout = setTimeout(() => {
        this.refresh()
      }, delay * 1000)
    }
  }

  get params(): Params {
    return this.paramsSubject.getValue();
  }

  get loading(): boolean {
    return this.loadingSubject.getValue()
  }

  setTimeframe(timeframe: TimeframeName) {
    this.updateParams((p) => ({
      ...p,
      buckets: {
        ...p.buckets,
        timeframe,
      }
    }))
  }

  setQueue(queue: string | null) {
    this.updateParams((p) => ({
      ...p,
      queue: queue ?? undefined,
    }))
  }

  setBucketDuration(duration: QueueMetricsBucketDuration, multiplier = 1) {
    this.updateParams((p) => ({
      ...p,
      buckets: {
        ...p.buckets,
        duration,
        multiplier,
      }
    }))
  }

  setAutoRefreshInterval(autoRefreshInterval: AutoRefreshInterval) {
    this.updateParams((p) => ({
      ...p,
      autoRefresh: autoRefreshInterval,
    }))
  }

  private updateParams(fn: (p: Params) => Params) {
    this.paramsSubject.next(fn(this.params));
  }

  refresh() {
    return this.request(this.variablesSubject.getValue())
  }

  private request(variables:  generated.QueueMetricsQueryVariables) {
    clearTimeout(this.refreshTimeout)
    this.loadingSubject.next(true)
    return this.apollo.query<generated.QueueMetricsQuery, generated.QueueMetricsQueryVariables>({
      query: generated.QueueMetricsDocument,
      variables,
      fetchPolicy: "no-cache",
    }).pipe(
      map((r) => {
        this.loadingSubject.next(false)
        this.rawResultSubject.next(r.data)
      })
    ).subscribe()
  }
}

const createVaraibles = (params: Params): generated.QueueMetricsQueryVariables=> ({
  input: {
    bucketDuration: params.buckets.duration,
    queues: params.queue ? [params.queue] : undefined,
    startTime: params.buckets.timeframe === "all" ? undefined : new Date(new Date().getTime() - (1000 * timeframeLengths[params.buckets.timeframe])).toISOString()
  }
})

const fromEntries = <K extends string, V>(entries: Array<[K, V]>): Partial<Record<K, V>>  => Object.fromEntries(entries) as Partial<Record<K, V>>

const createResult = (params: Params, rawResult: generated.QueueMetricsQuery): Result => {
  const {bucketParams,earliestBucket, latestBucket} = createBucketParams(params, rawResult)
  return {
    bucketParams,
    queues: Object
      .entries(rawResult.queue.metrics.reduce<Record<string, [StatusCounts, Partial<Record<EventName, EventBucketEntries>>]>>(
        (acc, next) => {
          if (next.queue !== (params.queue ?? next.queue)) {
            return acc
          }
          let createdAt: NormalizedBucket | undefined = normalizeBucket(next.createdAtBucket, bucketParams)
          if (earliestBucket && earliestBucket.index > createdAt.index) {
            createdAt = undefined;
          }
          let ranAt = next.ranAtBucket ? normalizeBucket(next.ranAtBucket, bucketParams) : undefined
          if (ranAt && (latestBucket.index < ranAt.index || (earliestBucket && earliestBucket.index > ranAt.index))) {
            ranAt = undefined;
          }
          if (next.queue !== params.queue && !createdAt && (
            !ranAt || next.status === "pending"
          )) {
            return acc
          }
          const [currentStatusCounts, currentEventBuckets] = acc[next.queue] ?? [
            emptyStatusCounts,
            []
          ]
          const currentLatency = next.latency ? toSeconds(parseDuration(next.latency)) : undefined
          return {
            ...acc,
            [next.queue]: [(next.status === "pending" ? createdAt : ranAt) ? {
              ...currentStatusCounts,
              [next.status]: next.count + currentStatusCounts[next.status],
            } : currentStatusCounts, {
              created: createdAt ? {
                ...currentEventBuckets.created,
                [createdAt.key]: {
                  count: next.count + (currentEventBuckets.created?.[createdAt.key]?.count ?? 0),
                  totalLatency: 0,
                  startTime: createdAt.start,
                }
              } : currentEventBuckets.created,
              processed: (ranAt && next.status === "processed") ? {
                ...currentEventBuckets.processed,
                [ranAt.key]: {
                  count: next.count + (currentEventBuckets.processed?.[ranAt.key]?.count ?? 0),
                  totalLatency: (currentEventBuckets.processed?.[ranAt.key]?.totalLatency ?? 0) + (currentLatency ?? 0),
                  startTime: ranAt.start,
                }
              } : currentEventBuckets.processed,
              failed: (ranAt && next.status === "failed") ? {
                ...currentEventBuckets.failed,
                [ranAt.key]: {
                  count: next.count + (currentEventBuckets.failed?.[ranAt.key]?.count ?? 0),
                  totalLatency: (currentEventBuckets.failed?.[ranAt.key]?.totalLatency ?? 0) + (currentLatency ?? 0),
                  startTime: ranAt.start,
                }
              } : currentEventBuckets.failed
            }]
          };
        },
        {}
      )).map(([queue, [statusCounts, eventBuckets]]) => {
        let events: QueueEvents | undefined;
        // const bucketKeys = Object.keys(eventBuckets).sort()
        if (Object.keys(eventBuckets).length) {
          const bucketDates = Array<number>()
          const buckets: EventBuckets = fromEntries(Array<EventName>("created", "processed", "failed").flatMap<[EventName, EventBucket]>((event): [EventName, EventBucket][] => {
            const entries = fromEntries(
              Object.entries(eventBuckets[event] ?? {}).filter(([, v]) => v?.count).sort(([a], [b]) => parseInt(a) < parseInt(b) ? 1 : -1),
            )
            const keys = Object.keys(entries)
            if (!keys.length) {
              return [];
            }
            const earliestBucket = parseInt(keys[0]);
            const latestBucket = parseInt(keys[keys.length - 1]);
            bucketDates.push(earliestBucket, latestBucket)
            return [[event, {
              earliestBucket,
              latestBucket,
              entries
            }]]
          }));
          bucketDates.sort()
          events = {
            bucketDuration: bucketParams.duration,
            earliestBucket: bucketDates[0],
            latestBucket: bucketDates[bucketDates.length - 1],
            eventBuckets: buckets,
          }
        }
        return {
          queue,
          statusCounts,
          events,
          isEmpty: !events?.eventBuckets
        }
      })}};

const createBucketParams = (params: Params, rawResult: generated.QueueMetricsQuery): {
  bucketParams: BucketParams<false>,
  earliestBucket: NormalizedBucket | undefined,
  latestBucket: NormalizedBucket
} => {
  const duration = params.buckets.duration
  let multiplier = params.buckets.multiplier === "AUTO" ? 1 : params.buckets.multiplier
  const timeframe= params.buckets.timeframe
  const now = new Date()
  const nowBucket = normalizeBucket(now, { duration, multiplier })
  const startBucket = timeframe === "all" ? undefined : normalizeBucket(now.getTime() - (1000 * timeframeLengths[timeframe]), {duration, multiplier})
  const allBuckets = [
    ...startBucket ? [startBucket] : [],
    ...rawResult.queue.metrics.flatMap((b) => [
    normalizeBucket(b.createdAtBucket, {duration, multiplier}),
    ...b.ranAtBucket ? [normalizeBucket(b.ranAtBucket, {duration, multiplier})] : [],
  ]),
    nowBucket
  ].filter((b) => !startBucket || b.index >= startBucket.index).sort((a, b) => a.index - b.index )
  const minBucket = allBuckets[0];
  const maxBucket = allBuckets[allBuckets.length - 1];
  if (params.buckets.multiplier === "AUTO") {
    const targetSpan = 20;
    const span = maxBucket.index - minBucket.index;
    multiplier = span === 0 ? 1 : Math.min(60, Math.max(Math.ceil(span / targetSpan), 1))
  }
  return {
    bucketParams: {
      duration,
      multiplier,
      timeframe,
    },
    earliestBucket: timeframe === "all" ? undefined : normalizeBucket(now.getTime() - (1000 * timeframeLengths[timeframe]), {duration, multiplier}),
    latestBucket: normalizeBucket(Math.max(now.getTime(), maxBucket.start.getTime()), { duration, multiplier })
  }
}

type NormalizedBucket = {
  key: string;
  index: number;
  start: Date;
}

export const normalizeBucket = (rawDate: string | number | Date, params: Pick<BucketParams<false>, "duration" | "multiplier">): NormalizedBucket => {
  const date = new Date(rawDate)
  const msMultiplier = 1000 * durationSeconds[params.duration] * params.multiplier
  const baseNumber = Math.floor(date.getTime() / msMultiplier)
  return {
    key: `${baseNumber}`,
    index: baseNumber,
    start: new Date(baseNumber * msMultiplier),
  }
}
