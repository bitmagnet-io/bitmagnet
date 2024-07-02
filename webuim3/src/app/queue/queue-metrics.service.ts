import {Apollo} from "apollo-angular";
import * as generated from '../graphql/generated';
import {map} from "rxjs/operators";
import {parse as parseDuration, toSeconds} from "iso8601-duration";
import {inject, Injectable} from "@angular/core";

export type Params = {
  bucketDuration: generated.QueueMetricsBucketDuration
}

type StatusCounts = Record<generated.QueueJobStatus, number>

type Event = "created" | "processed" | "failed"

type EventBucket = {
  count: number,
  latency?: number,
}

type QueueSummary = {
  queue: string
  statusCounts: StatusCounts
  eventBuckets: Record<Event, Record<string, EventBucket>>
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
        Object.values(r.data.queue.metrics.reduce<Record<string, QueueSummary>>(
          (acc, next) => {
            const current = acc[next.queue] ?? {
              queue: next.queue,
              statusCounts: {
                pending: 0,
                failed: 0,
                retry: 0,
                processed: 0,
              },
              eventBuckets: {
                created: {},
                processed: {},
                failed: {},
              }
            }
            const currentBucket: Record<Event, EventBucket> = {
              created: {
                count: current.eventBuckets.created[next.createdAtBucket]?.count ?? 0,
              },
              processed: next.ranAtBucket ? {
                count: current.eventBuckets.processed[next.ranAtBucket]?.count ?? 0,
                latency: current.eventBuckets.processed[next.ranAtBucket]?.latency,
              } : { count: 0 },
              failed: next.ranAtBucket ? {
                count: current.eventBuckets.failed[next.ranAtBucket]?.count ?? 0,
                latency: current.eventBuckets.failed[next.ranAtBucket]?.latency,
              } : { count: 0 },
            }
            const nextLatency = next.latency ? toSeconds(parseDuration(next.latency)) : undefined
            return {
              ...acc,
              [next.queue]: {
                queue: next.queue,
                statusCounts: {
                  ...current.statusCounts,
                  [next.status]: next.count + current.statusCounts[next.status],
                },
                eventBuckets: {
                  created: {
                    ...current.eventBuckets.created,
                    [next.createdAtBucket]: {
                      count: next.count + currentBucket.created.count,
                    }
                  },
                  processed: (next.ranAtBucket && next.status === "processed") ? {
                    ...current.eventBuckets.processed,
                    [next.ranAtBucket]: {
                      count: next.count + currentBucket.processed.count,
                      latency: (((nextLatency ?? 0) * next.count) + ((currentBucket.processed.latency ?? 0) * currentBucket.processed.count)) / (next.count + currentBucket.processed.count)
                    },
                  } : current.eventBuckets.processed,
                  failed: (next.ranAtBucket && next.status === "failed") ? {
                    ...current.eventBuckets.failed,
                    [next.ranAtBucket]: {
                      count: next.count + currentBucket.failed.count,
                      latency: (((nextLatency ?? 0) * next.count) + ((currentBucket.failed.latency ?? 0) * currentBucket.failed.count)) / (next.count + currentBucket.failed.count)
                    },
                  } : current.eventBuckets.failed
                }
              } as QueueSummary
            }
          },
          {}
        ))
      )
    )
  }
}
