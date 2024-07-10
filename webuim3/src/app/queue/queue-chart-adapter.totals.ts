import {ChartAdapter} from "../charting/types";
import {ChartConfiguration} from "chart.js";
import {Result} from "./queue-metrics.types";
import * as generated from "../graphql/generated";
import {statusNames} from "./queue.constants";
import {ThemeBaseColor} from "../themes/theme-types";
import {QueueJobStatus} from "../graphql/generated";
import {createThemeColor} from "../themes/theme-utils";
import {inject, Injectable} from "@angular/core";
import {ThemeInfoService} from "../themes/theme-info.service";
import {TranslocoService} from "@jsverse/transloco";

const statusColors: Record<QueueJobStatus, ThemeBaseColor> = {
  'pending': 'primary',
  'processed': 'success',
  'failed': 'error',
  'retry': 'caution'
}

@Injectable({providedIn: "root"})
export class QueueChartAdapterTotals implements ChartAdapter<Result> {
  private themeInfo = inject(ThemeInfoService)
  private transloco = inject(TranslocoService)

  create(result?: Result) : ChartConfiguration<"bar"> {
    const { colors } = this.themeInfo.info
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
          x: {
            ticks: {
              callback: (v) => parseInt(v as string).toLocaleString(this.transloco.getActiveLang())
            }
            },
          y: {
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
