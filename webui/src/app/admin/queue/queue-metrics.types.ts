import * as generated from "../../graphql/generated";
import {
  autoRefreshIntervalNames,
  eventNames,
  timeframeNames,
} from "./queue.constants";

export type TimeframeName = (typeof timeframeNames)[number];

export type BucketParams<_withAuto extends boolean = boolean> = {
  duration: _withAuto extends true
    ? "AUTO" | generated.MetricsBucketDuration
    : generated.MetricsBucketDuration;
  multiplier: _withAuto extends true ? "AUTO" | number : number;
  timeframe: TimeframeName;
};

export type AutoRefreshInterval = (typeof autoRefreshIntervalNames)[number];

export type StatusCounts = Record<generated.QueueJobStatus, number>;

export type EventName = (typeof eventNames)[number];

export type Params<_withAuto extends boolean = boolean> = {
  buckets: BucketParams<_withAuto>;
  queue?: string;
  event?: EventName;
  autoRefresh: AutoRefreshInterval;
};

export type EventBucketEntry = {
  startTime: Date;
  count: number;
  latency: number;
};

export type EventBucketEntries = Partial<Record<string, EventBucketEntry>>;

export type BucketSpan = {
  earliestBucket: number;
  latestBucket: number;
};

export type EventBucket = BucketSpan & {
  entries: EventBucketEntries;
};

export type EventBuckets = Partial<Record<EventName, EventBucket>>;

export type QueueEvents = BucketSpan & {
  bucketDuration: generated.MetricsBucketDuration;
  eventBuckets: EventBuckets;
};

export type QueueSummary<IsEmpty extends boolean = boolean> = {
  queue: string;
  isEmpty: IsEmpty;
  statusCounts: StatusCounts;
  events: IsEmpty extends false ? QueueEvents : undefined;
};

export type Result = {
  params: Params<false>;
  // earliestBucket: number;
  // latestBucket: number;
  bucketSpan?: BucketSpan;
  queues: QueueSummary[];
};
