import {ChartAdapter} from "../charting/types";
import {ChartConfiguration} from "chart.js";
import {BucketParams, Result} from "./queue-metrics.types";
import {durationSeconds, eventNames, timeframeLengths} from "./queue-metrics.constants";
import {normalizeBucket} from "./queue-metrics.controller";

const eventColors = {
  'created': 'blue',
  'processed': 'green',
  'failed': 'red',
}

export const queueChartAdapterTimeline: ChartAdapter<Result> = {
  create: (result) => {
    console.log({result})
    const labels = Array<string>()
    const datasets: ChartConfiguration<"line">["data"]["datasets"] = []
    if (result) {
      const nonEmptyQueues = result.queues.filter((q) => !q.isEmpty)
      const nonEmptyBuckets = Array.from(new Set(nonEmptyQueues.flatMap(
        (q) => q.events ? [q.events.earliestBucket, q.events.latestBucket] : []
      ))).sort()
      const now = new Date()
      const minBucket = result.bucketParams.timeframe ==="all"?nonEmptyBuckets[0]:Math.min(
        nonEmptyBuckets[0],
        normalizeBucket(now.getTime() - (1000 * timeframeLengths[result.bucketParams.timeframe]), result.bucketParams).index
      )
      const maxBucket = Math.max(
        nonEmptyBuckets[nonEmptyBuckets.length -1],
        normalizeBucket(now, result.bucketParams).index
      )
      // const seriesLabels = nonEmptyQueues.flatMap((q) => events.map((status) => [q.queue, status].join("/")))
      if (nonEmptyBuckets.length) {
        for (let i = minBucket; i <= maxBucket; i++) {
          labels.push(formatBucketKey(result.bucketParams, i))
        }
        for (const queue of nonEmptyQueues) {
          for (const event of eventNames) {
            const series = Array<number>()
            for (let i = minBucket; i <= maxBucket; i++) {
              series.push(queue.events?.eventBuckets?.[event]?.entries?.[`${i}`]?.count ?? 0)
            }
            datasets.push({
              label: [queue.queue, event].join("/"),
              data: series,
              // fill: 'origin',
              // backgroundColor: 'rgba(148,159,177,0.2)',
              borderColor: eventColors[event],
              pointBackgroundColor: 'rgba(148,159,177,1)',
              pointBorderColor: '#fff',
              pointHoverBackgroundColor: '#fff',
              pointHoverBorderColor: 'rgba(148,159,177,0.8)',
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
          y: {
            position: 'left',
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
