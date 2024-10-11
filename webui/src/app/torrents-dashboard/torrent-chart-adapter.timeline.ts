import { ChartConfiguration } from "chart.js";
import { inject, Injectable } from "@angular/core";
import { TranslocoService } from "@jsverse/transloco";
import { format as formatDate } from "date-fns/format";
import { ChartAdapter } from "../charting/types";
import { ThemeBaseColor } from "../themes/theme-types";
import { createThemeColor } from "../themes/theme-utils";
import { ThemeInfoService } from "../themes/theme-info.service";
import { resolveDateLocale } from "../dates/dates.locales";
import { formatDuration } from "../dates/dates.utils";
import { normalizeBucket } from "./torrent-metrics.controller";
import {
  durationSeconds,
  eventNames,
  timeframeLengths,
} from "./torrent-metrics.constants";
import { BucketParams, EventName, Result } from "./torrent-metrics.types";

const eventColors: Record<EventName, ThemeBaseColor> = {
  created: "primary",
  updated: "secondary",
};

@Injectable({ providedIn: "root" })
export class TorrentChartAdapterTimeline
  implements ChartAdapter<Result, "line">
{
  private themeInfo = inject(ThemeInfoService);
  private transloco = inject(TranslocoService);

  create(result?: Result): ChartConfiguration<"line"> {
    const { colors } = this.themeInfo.info;
    const labels = Array<string>();
    const datasets: ChartConfiguration<"line">["data"]["datasets"] = [];
    if (result) {
      const nonEmptySources = result.sources.filter((q) => !q.isEmpty);
      const nonEmptyBuckets = Array.from(
        new Set(
          nonEmptySources.flatMap((q) =>
            q.events ? [q.events.earliestBucket, q.events.latestBucket] : [],
          ),
        ),
      ).sort();
      const now = new Date();
      const minBucket = Math.min(
        nonEmptyBuckets[0],
        normalizeBucket(
          now.getTime() -
            1000 * timeframeLengths[result.params.buckets.timeframe],
          result.params.buckets,
        ).index,
      );
      const maxBucket = Math.max(
        nonEmptyBuckets[nonEmptyBuckets.length - 1],
        normalizeBucket(now, result.params.buckets).index,
      );
      if (nonEmptyBuckets.length) {
        for (let i = minBucket; i <= maxBucket; i++) {
          labels.push(this.formatBucketKey(result.params.buckets, i));
        }
        const relevantEvents = eventNames.filter(
          (n) => (result.params.event ?? n) === n,
        );
        for (const source of nonEmptySources) {
          for (const event of relevantEvents) {
            const series = Array<number>();
            for (let i = minBucket; i <= maxBucket; i++) {
              series.push(
                source.events?.eventBuckets?.[event]?.entries?.[`${i}`]
                  ?.count ?? 0,
              );
            }
            datasets.push({
              yAxisID: "yCount",
              label: [source.source, event].join("/"),
              data: series,
              // fill: 'origin',
              // backgroundColor: 'rgba(148,159,177,0.2)',
              borderColor: colors[createThemeColor(eventColors[event], 50)],
              pointBackgroundColor:
                colors[createThemeColor(eventColors[event], 20)],
              pointBorderColor:
                colors[createThemeColor(eventColors[event], 80)],
              pointHoverBackgroundColor:
                colors[createThemeColor(eventColors[event], 40)],
              pointHoverBorderColor:
                colors[createThemeColor(eventColors[event], 60)],
            });
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
            position: "left",
            ticks: {
              callback: (v) =>
                parseInt(v as string).toLocaleString(
                  this.transloco.getActiveLang(),
                ),
            },
          },
          // x: {
          //   ticks: {
          //     stepSize: 5
          //   }
          // }
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
    };
  }

  private formatBucketKey(params: BucketParams<false>, key: number): string {
    let formatStr: string;
    switch (params.duration) {
      case "day":
        formatStr = "d LLL";
        break;
      case "hour":
        formatStr = "d LLL H:00";
        break;
      case "minute":
        formatStr = "H:mm";
        break;
    }
    return formatDate(
      1000 * durationSeconds[params.duration] * params.multiplier * key,
      formatStr,
      {
        locale: resolveDateLocale(this.transloco.getActiveLang()),
      },
    );
  }

  private formatDuration(d: number | string): string {
    if (typeof d === "string") {
      d = parseInt(d);
    }
    if (d === 0) {
      return "0";
    }
    let seconds = d;
    let minutes = 0;
    let hours = 0;
    let days = 0;
    if (seconds >= 60) {
      minutes = Math.floor(seconds / 60);
      seconds = seconds % 60;
      if (minutes >= 5) {
        seconds = 0;
        if (minutes >= 60) {
          hours = Math.floor(minutes / 60);
          minutes = minutes % 60;
          if (hours >= 5) {
            minutes = 0;
            if (hours >= 24) {
              days = Math.floor(hours / 24);
              hours = hours % 24;
            }
          }
        }
      }
    }
    return formatDuration(
      { days, hours, minutes, seconds },
      this.transloco.getActiveLang(),
    );
  }
}
