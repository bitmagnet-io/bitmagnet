import * as generated from "../../graphql/generated";
import {
  autoRefreshIntervalNames,
  eventNames,
  timeframeNames,
} from "./torrent-metrics.constants";

export type TimeframeName = (typeof timeframeNames)[number];

export type BucketParams<_withAuto extends boolean = boolean> = {
  duration: _withAuto extends true
    ? "AUTO" | generated.MetricsBucketDuration
    : generated.MetricsBucketDuration;
  multiplier: _withAuto extends true ? "AUTO" | number : number;
  timeframe: TimeframeName;
};

export type AutoRefreshInterval = (typeof autoRefreshIntervalNames)[number];

export type EventName = (typeof eventNames)[number];

export type Params<_withAuto extends boolean = boolean> = {
  buckets: BucketParams<_withAuto>;
  source?: string;
  event?: EventName;
  autoRefresh: AutoRefreshInterval;
};

export type EventBucketEntry = {
  startTime: Date;
  count: number;
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

export type TorrentEvents = BucketSpan & {
  bucketDuration: generated.MetricsBucketDuration;
  eventBuckets: EventBuckets;
};

export type SourceSummary<IsEmpty extends boolean = boolean> = {
  source: string;
  isEmpty: IsEmpty;
  events: IsEmpty extends false ? TorrentEvents : undefined;
};

export type Result = {
  params: Params<false>;
  bucketSpan?: BucketSpan;
  sourceSummaries: SourceSummary[];
  availableSources: generated.TorrentSource[];
};
