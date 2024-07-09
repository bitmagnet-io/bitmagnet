import {ChartAdapter} from "../charting/types";
import {ChartConfiguration} from "chart.js";
import {BucketParams, EventName, Result} from "./queue-metrics.types";
import {durationSeconds, eventNames, timeframeLengths} from "./queue.constants";
import {normalizeBucket} from "./queue-metrics.controller";
import {ThemeBaseColor, ThemeColor} from "../themes/theme-types";
import {createThemeColor} from "../themes/theme-utils";

const eventColors: Record<EventName, ThemeBaseColor> = {
  'created': 'primary',
  'processed': 'success',
  'failed': 'error',
}

export const queueChartAdapterTimeline: ChartAdapter<Result> = {
  create: (result, {colors}) => {
    const labels = Array<string>()
    const datasets: ChartConfiguration<"line">["data"]["datasets"] = []
    if (result) {
      const nonEmptyQueues = result.queues.filter((q) => !q.isEmpty)
      const nonEmptyBuckets = Array.from(new Set(nonEmptyQueues.flatMap(
        (q) => q.events ? [q.events.earliestBucket, q.events.latestBucket] : []
      ))).sort()
      const now = new Date()
      const minBucket = result.params.buckets.timeframe ==="all"?nonEmptyBuckets[0]:Math.min(
        nonEmptyBuckets[0],
        normalizeBucket(now.getTime() - (1000 * timeframeLengths[result.params.buckets.timeframe]), result.params.buckets).index
      )
      const maxBucket = Math.max(
        nonEmptyBuckets[nonEmptyBuckets.length -1],
        normalizeBucket(now, result.params.buckets).index
      )
      // const seriesLabels = nonEmptyQueues.flatMap((q) => events.map((status) => [q.queue, status].join("/")))
      if (nonEmptyBuckets.length) {
        for (let i = minBucket; i <= maxBucket; i++) {
          labels.push(formatBucketKey(result.params.buckets, i))
        }
        const relevantEvents = eventNames.filter((n) => (result.params.event ?? n) === n)
        for (const queue of nonEmptyQueues) {
          for (const event of relevantEvents) {
            const series = Array<number>()
            for (let i = minBucket; i <= maxBucket; i++) {
              series.push(queue.events?.eventBuckets?.[event]?.entries?.[`${i}`]?.count ?? 0)
            }
            datasets.push({
              yAxisID: "yCount",
              label: [queue.queue, event].join("/"),
              data: series,
              // fill: 'origin',
              // backgroundColor: 'rgba(148,159,177,0.2)',
              borderColor: colors[createThemeColor(eventColors[event], 50)],
              pointBackgroundColor: colors[createThemeColor(eventColors[event], 20)],
              pointBorderColor: colors[createThemeColor(eventColors[event], 80)],
              pointHoverBackgroundColor: colors[createThemeColor(eventColors[event], 40)],
              pointHoverBorderColor: colors[createThemeColor(eventColors[event], 60)],
            })
          }
          const latencyEvents = (["processed", "failed"] as const).filter((e) => relevantEvents.includes(e))
          if (latencyEvents.length) {
            const latencySeries = Array<number|null>()
            for (let i = minBucket; i <= maxBucket; i++) {
              const result = (["processed", "failed"] as const).filter((e) => relevantEvents.includes(e)).reduce<[number, number] | null>(
                (acc, next) => {
                  const entry = queue.events?.eventBuckets?.[next]?.entries?.[`${i}`]
                  if (!entry?.count) {
                    return acc
                  }
                  return [(acc?.[0] ?? 0)+entry.latency, (acc?.[1] ?? 0)+entry.count]
                },
                null
              )
              latencySeries.push(result ? result[0] / result[1] : null)
            }
            datasets.push({
              yAxisID: "yLatency",
              label: [queue.queue, 'latency'].join("/"),
              data: latencySeries,
              // fill: 'origin',
              // backgroundColor: 'rgba(148,159,177,0.2)',
              borderColor: colors["tertiary-50"],
              // pointBackgroundColor: 'rgba(148,159,177,1)',
              // pointBorderColor: '#fff',
              pointHoverBackgroundColor: colors["tertiary-80"],
              pointHoverBorderColor: colors["tertiary-20"],
            })
          }
        }
      }
    }
    return {
      type: "line",
      options: {
        animation: false,
        elements: {
          line: {
            tension: 0.5,
          },
        },
        scales: {
          yCount: {
            position: 'left',
            // max: 100,
            ticks: {
              stepSize: 1
            }
          },
          yLatency: {
            position: 'right',
            // max: 100,
          },
          x: {
            ticks: {
              stepSize: 5
            }
          }
          // y1: {
          //   position: 'right',
          //   grid: {
          //     color: 'rgba(255,0,0,0.3)',
          //   },
          //   ticks: {
          //     color: 'red',
          //   },
          // },
        },
        plugins: {
          legend: {
            display: true,
          },
          decimation: {
            enabled: true,
          }
          // datalabels: {
          //   anchor: 'end',
          //   align: 'end',
          // },
        },
      },
      data: {
        labels,
        datasets,
      },
    }
  }
}

const formatBucketKey = (params: BucketParams<false>, key: number) => {
  const msMultiplier = 1000 * durationSeconds[params.duration] * params.multiplier
  return new Date(key * msMultiplier).toISOString().split("T")[1]
}
