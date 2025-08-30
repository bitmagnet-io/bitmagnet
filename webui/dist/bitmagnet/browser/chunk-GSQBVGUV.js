// src/app/dashboard/queue/queue.constants.ts
var emptyStatusCounts = {
  pending: 0,
  failed: 0,
  retry: 0,
  processed: 0
};
var defaultBucketParams = {
  duration: "hour",
  multiplier: 1,
  timeframe: "all"
};
var resolutionNames = ["day", "hour", "minute"];
var durationSeconds = {
  minute: 60,
  hour: 60 * 60,
  day: 60 * 60 * 24
};
var emptyParams = {
  buckets: defaultBucketParams,
  autoRefresh: "off"
};
var emptyResult = {
  params: emptyParams,
  queues: []
};
var eventNames = ["created", "processed", "failed"];
var statusNames = ["pending", "processed", "retry", "failed"];
var timeframeNames = [
  "minutes_15",
  "minutes_30",
  "hours_1",
  "hours_6",
  "hours_12",
  "days_1",
  "weeks_1",
  "all"
];
var timeframeLengths = {
  minutes_15: 60 * 15,
  minutes_30: 60 * 30,
  hours_1: 60 * 60,
  hours_6: 60 * 60 * 6,
  hours_12: 60 * 60 * 12,
  days_1: 60 * 60 * 24,
  weeks_1: 60 * 60 * 24 * 7,
  all: Infinity
};
var availableQueueNames = ["process_torrent", "process_torrent_batch"];
var autoRefreshIntervalNames = [
  "off",
  "seconds_10",
  "seconds_30",
  "minutes_1",
  "minutes_5"
];
var autoRefreshIntervals = {
  off: null,
  seconds_10: 10,
  seconds_30: 30,
  minutes_1: 60,
  minutes_5: 60 * 5
};

export {
  emptyStatusCounts,
  resolutionNames,
  durationSeconds,
  emptyParams,
  emptyResult,
  eventNames,
  statusNames,
  timeframeNames,
  timeframeLengths,
  availableQueueNames,
  autoRefreshIntervalNames,
  autoRefreshIntervals
};
//# sourceMappingURL=chunk-GSQBVGUV.js.map
