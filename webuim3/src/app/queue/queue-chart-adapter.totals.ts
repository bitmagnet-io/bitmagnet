import {ChartAdapter} from "../charting/types";
import {ChartConfiguration} from "chart.js";
import {Result} from "./queue-metrics.types";

export const queueChartAdapterTotals: ChartAdapter<Result> = {
  create: (result) => {
    const labels = Array<string>()
    const datasets: ChartConfiguration<"bar">["data"]["datasets"] = []
    if (result) {
      const bucketKeys = Array.from(new Set(result.queues.flatMap((q) => q.events ? [q.events.earliestBucket, q.events.latestBucket] : []))).sort()
      if (bucketKeys.length) {
        const nonEmptyQueues = result.queues.filter((q) => !q.isEmpty)
        labels.push(...nonEmptyQueues.map((q) => q.queue))
        datasets.push(...(["pending", "retry", "failed", "processed"] as const).map((status) => ({
          label: status,
          data: nonEmptyQueues.map((q) => q.statusCounts[status]),
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
