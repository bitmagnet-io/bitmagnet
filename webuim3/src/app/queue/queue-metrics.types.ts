import * as generated from "../graphql/generated";
import {autoRefreshIntervalNames, autoRefreshIntervals, eventNames, timeframeNames} from "./queue-metrics.constants";

export type TimeframeName = typeof timeframeNames[number]

export type BucketParams<_withAuto extends boolean = boolean> = {
  duration: generated.QueueMetricsBucketDuration
  multiplier: _withAuto extends true ? "AUTO" | number : number;
  timeframe: TimeframeName;
}

export type AutoRefreshInterval = typeof autoRefreshIntervalNames[number]

export type Params = {
  buckets: BucketParams<true>;
  queue?: string;
  autoRefresh: AutoRefreshInterval;
}

export type StatusCounts = Record<generated.QueueJobStatus, number>

export type EventName = typeof eventNames[number]

export type EventBucketEntry = {
  startTime: Date;
  count: number,
  totalLatency: number,
}

export type EventBucketEntries = Partial<Record<string, EventBucketEntry>>

type BucketSpan = {
  earliestBucket: number;
  latestBucket: number;
}

export type EventBucket = BucketSpan & {
  entries: EventBucketEntries
};

export type EventBuckets = Partial<Record<EventName, EventBucket>>

export type QueueEvents = BucketSpan & {
  bucketDuration: generated.QueueMetricsBucketDuration
  eventBuckets: EventBuckets;
}

export type QueueSummary<IsEmpty extends boolean = boolean> = {
  queue: string
  isEmpty: IsEmpty
  statusCounts: StatusCounts
  events: IsEmpty extends false ? QueueEvents : undefined;
}

export type Result = {
  bucketParams: BucketParams<false>;
  // earliestBucket: number;
  // latestBucket: number;
  bucketSpan?: BucketSpan;
  queues: QueueSummary[]
}
