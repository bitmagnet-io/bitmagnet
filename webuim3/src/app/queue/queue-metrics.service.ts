import {Apollo} from "apollo-angular";
import * as generated from '../graphql/generated';
import {map} from "rxjs/operators";
import {parse as parseDuration, toSeconds} from "iso8601-duration";
import {inject, Injectable} from "@angular/core";

export type Params = {
  bucketDuration: generated.QueueMetricsBucketDuration
}

type StatusCounts = Record<generated.QueueJobStatus, number>

const emptyStatusCounts: StatusCounts = {
  pending: 0,
  failed: 0,
  retry: 0,
  processed: 0,
}

type Event = "created" | "processed" | "failed"

type EventBucketEntry = {
  count: number,
  totalLatency: number,
}

type EventBucketEntries = Partial<Record<string, EventBucketEntry>>

type EventBucket = {
  earliestBucket: string;
  latestBucket: string;
  entries: EventBucketEntries
};

type EventBuckets = Partial<Record<Event, EventBucket>>

type QueueEvents = {
  bucketDuration: generated.QueueMetricsBucketDuration
  earliestBucket: string;
  latestBucket: string;
  eventBuckets: EventBuckets;
}

type QueueSummary = {
  queue: string
  statusCounts: StatusCounts
  events?: QueueEvents
}

export type Result = {
  params: Params
  queues: QueueSummary[]
}

@Injectable({providedIn: 'root'})
export class QueueMetricsService {
  private apollo = inject(Apollo)

  request(params: Params) {
    return this.apollo.query<generated.QueueMetricsQuery, generated.QueueMetricsQueryVariables>({
      query: generated.QueueMetricsDocument,
      variables: {
        input: {
          bucketDuration: params.bucketDuration,
        }
      }
    }).pipe(
      map((r): QueueSummary[] =>
        Object.entries(r.data.queue.metrics.reduce<Record<string, [StatusCounts, Partial<Record<Event, EventBucketEntries>>]>>(
          (acc, next) => {
            const [currentStatusCounts, currentEventBuckets ] = acc[next.queue] ?? [
              emptyStatusCounts,
              []
            ]
            const currentLatency = next.latency ? toSeconds(parseDuration(next.latency)) : undefined
            const currentBucket: Record<Event, EventBucketEntry> = {
              created: {
                count: next.count + (currentEventBuckets.created?.[next.createdAtBucket]?.count ?? 0),
                totalLatency: 0,
              },
              processed: next.ranAtBucket ? {
                count: next.count + (currentEventBuckets.processed?.[next.ranAtBucket]?.count ?? 0),
                totalLatency: (currentEventBuckets.processed?.[next.ranAtBucket]?.totalLatency ?? 0) + (currentLatency ?? 0),
              } : { count: 0, totalLatency: 0 },
              failed: next.ranAtBucket ? {
                count: next.count + (currentEventBuckets.failed?.[next.ranAtBucket]?.count ?? 0),
                totalLatency: (currentEventBuckets.failed?.[next.ranAtBucket]?.totalLatency ?? 0) + (currentLatency ?? 0),
              } : { count: 0, totalLatency: 0 },
            }
            return {
              ...acc,
              [next.queue]: [{
                ...currentStatusCounts,
                [next.status]: next.count + currentStatusCounts[next.status],
              },  {
                created: {
                  ...currentEventBuckets.created,
                  [next.createdAtBucket]: currentBucket.created
                },
                processed: (next.ranAtBucket && next.status === "processed") ? {
                  ...currentEventBuckets.processed,
                  [next.ranAtBucket]: currentBucket.processed
                } : currentEventBuckets.processed,
                failed: (next.ranAtBucket && next.status === "failed") ? {
                  ...currentEventBuckets.failed,
                  [next.ranAtBucket]: currentBucket.failed
                } : currentEventBuckets.failed
              }]
            };
          },
          {}
        )).map(([queue, [statusCounts, eventBuckets]]) => {
          let events: QueueEvents | undefined;
          // const bucketKeys = Object.keys(eventBuckets).sort()
          if (Object.keys(eventBuckets).length) {
            const bucketDates = Array<string>()
            const buckets: EventBuckets = fromEntries(Array<Event>("created", "processed", "failed").flatMap<[Event, EventBucket]>((event): [Event, EventBucket][] => {
              const entries = fromEntries(
                Object.entries(eventBuckets[event] ?? {}).filter(([, v]) => v?.count).sort(([a], [b]) => a < b ? 1 : -1),
              )
              const keys = Object.keys(entries)
              if (!keys.length) {
                return [];
              }
              const earliestBucket = keys[0];
              const latestBucket = keys[keys.length - 1];
              bucketDates.push(earliestBucket, latestBucket)
              return [[event, {
                earliestBucket,
                latestBucket,
                entries
              }]]
            }));
            bucketDates.sort()
            events = {
              bucketDuration: params.bucketDuration,
              earliestBucket: bucketDates[0],
              latestBucket: bucketDates[bucketDates.length - 1],
              eventBuckets: buckets,
            }
          }
          return {
            queue,
            statusCounts,
            events,
          }
        })
      )
    )
  }
}

const fromEntries = <K extends string, V>(entries: Array<[K, V]>): Partial<Record<K, V>>  => Object.fromEntries(entries) as Partial<Record<K, V>>
