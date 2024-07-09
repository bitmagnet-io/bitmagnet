import {ChartAdapter} from "../charting/types";
import {ChartConfiguration} from "chart.js";
import {Result} from "./queue-metrics.types";
import * as generated from "../graphql/generated";
import {statusNames} from "./queue.constants";
import {ThemeBaseColor} from "../themes/theme-types";
import {QueueJobStatus} from "../graphql/generated";
import {createThemeColor} from "../themes/theme-utils";

const statusColors: Record<QueueJobStatus, ThemeBaseColor> = {
  'pending': 'primary',
  'processed': 'success',
  'failed': 'error',
  'retry': 'caution'
}

export const queueChartAdapterTotals: ChartAdapter<Result> = {
  create: (result, {colors}) => {
    const labels = Array<string>()
    const datasets: ChartConfiguration<"bar">["data"]["datasets"] = []
    if (result) {
      const bucketKeys = Array.from(new Set(result.queues.flatMap((q) => q.events ? [q.events.earliestBucket, q.events.latestBucket] : []))).sort()
      if (bucketKeys.length) {
        const nonEmptyQueues = result.queues.filter((q) => !q.isEmpty)
        labels.push(...nonEmptyQueues.map((q) => q.queue))
        const statuses = Array<generated.QueueJobStatus>()
        switch (result.params.event) {
          case "created":
            statuses.push("pending")
            break
          case "processed":
            statuses.push("processed")
            break
          case "failed":
            statuses.push("retry", "failed")
            break
          default:
            statuses.push(...statusNames)
            break
        }
        datasets.push(...statuses.map((status) => ({
          label: status,
          data: nonEmptyQueues.map((q) => q.statusCounts[status]),
          backgroundColor: colors[createThemeColor(statusColors[status], 50)],
        })))
      }
    }
    return {
      type: "bar",
      options: {
        animation: false,
        // transitions: {
        //   show: {
        //     animation: {
        //       duration: 1
        //     }
        //   },
        //   // animation: {
        //   //   duration: 10,
        //   // },
        //   // show: {
        //   //   animations: {
        //   //     duration: 10,
        //   //     // x: {
        //   //     //   from: 0
        //   //     // },
        //   //     // y: {
        //   //     //   from: 0
        //   //     // }
        //   //   }
        //   // },
        // },
        scales: {
          x: {},
          y: {
            // min: 10,
          },
        },
        indexAxis: "y",
        plugins: {
          legend: {
            display: true,
          },
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
