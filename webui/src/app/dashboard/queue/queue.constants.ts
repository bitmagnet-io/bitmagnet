import * as generated from "../../graphql/generated";
import {
  BucketParams,
  Params,
  Result,
  StatusCounts,
} from "./queue-metrics.types";

export const emptyStatusCounts: StatusCounts = {
  pending: 0,
  failed: 0,
  retry: 0,
  processed: 0,
};

export const defaultBucketParams: BucketParams = {
  duration: "hour",
  multiplier: 1,
  timeframe: "all",
};

export const resolutionNames = ["day", "hour", "minute"] as const;

export const durationSeconds: Record<generated.MetricsBucketDuration, number> =
  {
    minute: 60,
    hour: 60 * 60,
    day: 60 * 60 * 24,
  };

export const emptyParams: Params = {
  buckets: defaultBucketParams,
  autoRefresh: "off",
};

export const emptyResult: Result = {
  params: emptyParams,
  queues: [],
};

export const eventNames = ["created", "processed", "failed"] as const;

export const statusNames = ["pending", "processed", "retry", "failed"] as const;

export const timeframeNames = [
  "minutes_15",
  "minutes_30",
  "hours_1",
  "hours_6",
  "hours_12",
  "days_1",
  "weeks_1",
  "all",
] as const;

export const timeframeLengths: Record<(typeof timeframeNames)[number], number> =
  {
    minutes_15: 60 * 15,
    minutes_30: 60 * 30,
    hours_1: 60 * 60,
    hours_6: 60 * 60 * 6,
    hours_12: 60 * 60 * 12,
    days_1: 60 * 60 * 24,
    weeks_1: 60 * 60 * 24 * 7,
    all: Infinity,
  };

export const availableQueueNames = ["process_torrent", "process_torrent_batch"];

export const autoRefreshIntervalNames = [
  "off",
  "seconds_10",
  "seconds_30",
  "minutes_1",
  "minutes_5",
] as const;

export const autoRefreshIntervals: Record<
  (typeof autoRefreshIntervalNames)[number],
  number | null
> = {
  off: null,
  seconds_10: 10,
  seconds_30: 30,
  minutes_1: 60,
  minutes_5: 60 * 5,
};
