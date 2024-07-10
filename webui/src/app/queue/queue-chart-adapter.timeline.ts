import { ChartConfiguration } from 'chart.js';
import { inject, Injectable } from '@angular/core';
import { TranslocoService } from '@jsverse/transloco';
import { ChartAdapter } from '../charting/types';
import { ThemeBaseColor } from '../themes/theme-types';
import { createThemeColor } from '../themes/theme-utils';
import { ThemeInfoService } from '../themes/theme-info.service';
import { normalizeBucket } from './queue-metrics.controller';
import {
  durationSeconds,
  eventNames,
  timeframeLengths,
} from './queue.constants';
import { BucketParams, EventName, Result } from './queue-metrics.types';
import {format as formatDate} from "date-fns/format";
import {resolveDateLocale} from "../dates/dates.locales";
import {formatDuration} from "../dates/dates.utils";

const eventColors: Record<EventName, ThemeBaseColor> = {
  created: 'primary',
  processed: 'success',
  failed: 'error',
};

@Injectable({ providedIn: 'root' })
export class QueueChartAdapterTimeline implements ChartAdapter<Result, 'line'> {
  private themeInfo = inject(ThemeInfoService);
  private transloco = inject(TranslocoService);

  create(result?: Result): ChartConfiguration<'line'> {
    const { colors } = this.themeInfo.info;
    const labels = Array<string>();
    const datasets: ChartConfiguration<'line'>['data']['datasets'] = [];
    if (result) {
      const nonEmptyQueues = result.queues.filter((q) => !q.isEmpty);
      const nonEmptyBuckets = Array.from(
        new Set(
          nonEmptyQueues.flatMap((q) =>
            q.events ? [q.events.earliestBucket, q.events.latestBucket] : [],
          ),
        ),
      ).sort();
      const now = new Date();
      const minBucket =
        result.params.buckets.timeframe === 'all'
          ? nonEmptyBuckets[0]
          : Math.min(
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
        for (const queue of nonEmptyQueues) {
          for (const event of relevantEvents) {
            const series = Array<number>();
            for (let i = minBucket; i <= maxBucket; i++) {
              series.push(
                queue.events?.eventBuckets?.[event]?.entries?.[`${i}`]?.count ??
                  0,
              );
            }
            datasets.push({
              yAxisID: 'yCount',
              label: [queue.queue, event].join('/'),
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
          const latencyEvents = (['processed', 'failed'] as const).filter((e) =>
            relevantEvents.includes(e),
          );
          if (latencyEvents.length) {
            const latencySeries = Array<number | null>();
            for (let i = minBucket; i <= maxBucket; i++) {
              const result = (['processed', 'failed'] as const)
                .filter((e) => relevantEvents.includes(e))
                .reduce<[number, number] | null>((acc, next) => {
                  const entry =
                    queue.events?.eventBuckets?.[next]?.entries?.[`${i}`];
                  if (!entry?.count) {
                    return acc;
                  }
                  return [
                    (acc?.[0] ?? 0) + entry.latency,
                    (acc?.[1] ?? 0) + entry.count,
                  ];
                }, null);
              latencySeries.push(result ? result[0] / result[1] : null);
            }
            datasets.push({
              yAxisID: 'yLatency',
              label: [queue.queue, 'latency'].join('/'),
              data: latencySeries,
              // fill: 'origin',
              // backgroundColor: 'rgba(148,159,177,0.2)',
              borderColor: colors['tertiary-50'],
              // pointBackgroundColor: 'rgba(148,159,177,1)',
              // pointBorderColor: '#fff',
              pointHoverBackgroundColor: colors['tertiary-80'],
              pointHoverBorderColor: colors['tertiary-20'],
            });
          }
        }
      }
    }
    return {
      type: 'line',
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
            ticks: {
              callback: (v) =>
                parseInt(v as string).toLocaleString(
                  this.transloco.getActiveLang(),
                ),
            },
          },
          yLatency: {
            position: 'right',
            ticks: {
              callback: (v) => {
                if (typeof v === 'string') {
                  v = parseInt(v);
                }
                if (v === 0) {
                  return '0';
                }
                return formatDuration(v, this.transloco.getActiveLang());
                // if (v > 60 * 60) {
                //   return d.format('H[h]mm');
                // }
                // if (v < 1) {
                //   return d.format('SSS[ms]');
                // }
                // if (v < 5) {
                //   return d.format('s.SSS[s]');
                // }
                // return d.format('mm:ss[m]');
              },
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
    let formatStr: string
    switch (params.duration) {
      case "day":
        formatStr = "d LLL"
        break
      case "hour":
        formatStr = "d LLL H:00"
        break
      case "minute":
        formatStr = "H:mm"
        break
    }
    return formatDate(1000 * durationSeconds[params.duration] * params.multiplier * key, formatStr, {
      locale: resolveDateLocale(this.transloco.getActiveLang())
    })
  }
}
