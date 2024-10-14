import { ChartConfiguration } from "chart.js";
import { inject, Injectable } from "@angular/core";
import { TranslocoService } from "@jsverse/transloco";
import { format as formatDate } from "date-fns/format";
import { ChartAdapter } from "../../charting/types";
import { ThemeBaseColor } from "../../themes/theme-types";
import { createThemeColor } from "../../themes/theme-utils";
import { ThemeInfoService } from "../../themes/theme-info.service";
import { resolveDateLocale } from "../../dates/dates.locales";
import {
  durationSeconds,
  eventNames,
  timeframeLengths,
} from "./torrent-metrics.constants";
import { BucketParams, EventName, Result } from "./torrent-metrics.types";
import { normalizeBucket } from "./torrent-metrics.utils";

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
      const nonEmptySources = result.sourceSummaries.filter((q) => !q.isEmpty);
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
        },
        plugins: {
          legend: {
            display: true,
          },
          decimation: {
            enabled: true,
          },
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
}
