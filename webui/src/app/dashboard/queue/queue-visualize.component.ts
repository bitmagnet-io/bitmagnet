import { Component, inject, OnDestroy, OnInit } from "@angular/core";
import { Apollo } from "apollo-angular";
import { GraphQLModule } from "../../graphql/graphql.module";
import { ChartComponent } from "../../charting/chart.component";
import { BreakpointsService } from "../../layout/breakpoints.service";
import { ErrorsService } from "../../errors/errors.service";
import { AppModule } from "../../app.module";
import {
  autoRefreshIntervalNames,
  availableQueueNames,
  eventNames,
  resolutionNames,
  timeframeNames,
} from "./queue.constants";
import { QueueModule } from "./queue.module";
import { QueueChartAdapterTotals } from "./queue-chart-adapter.totals";
import { QueueMetricsController } from "./queue-metrics.controller";
import { QueueChartAdapterTimeline } from "./queue-chart-adapter.timeline";

@Component({
  selector: "app-queue-visualize",
  standalone: true,
  templateUrl: "./queue-visualize.component.html",
  styleUrl: "./queue-visualize.component.scss",
  imports: [AppModule, ChartComponent, GraphQLModule, QueueModule],
})
export class QueueVisualizeComponent implements OnInit, OnDestroy {
  breakpoints = inject(BreakpointsService);
  private apollo = inject(Apollo);
  queueMetricsController = new QueueMetricsController(
    this.apollo,
    {
      buckets: {
        duration: "AUTO",
        multiplier: "AUTO",
        timeframe: "all",
      },
      autoRefresh: "seconds_30",
    },
    inject(ErrorsService),
  );
  protected readonly timeline = inject(QueueChartAdapterTimeline);
  protected readonly totals = inject(QueueChartAdapterTotals);

  protected readonly resolutionNames = resolutionNames;
  protected readonly timeframeNames = timeframeNames;
  protected readonly availableQueueNames = availableQueueNames;
  protected readonly autoRefreshIntervalNames = autoRefreshIntervalNames;

  ngOnInit() {
    this.queueMetricsController.result$.subscribe((result) => {
      // change the default settings to more sensible ones if there is <12 hours of data to show
      if (
        this.queueMetricsController.params.buckets.timeframe === "all" &&
        this.queueMetricsController.params.buckets.duration === "AUTO" &&
        result.params.buckets.duration === "hour"
      ) {
        const span = result.bucketSpan;
        if (span && span.latestBucket - span.earliestBucket < 12) {
          this.queueMetricsController.setBucketDuration("minute");
        }
      }
    });
  }

  ngOnDestroy() {
    this.queueMetricsController.setAutoRefreshInterval("off");
  }

  protected readonly eventNames = eventNames;

  handleMultiplierEvent(event: Event) {
    const value = (event.currentTarget as HTMLInputElement).value;
    this.queueMetricsController.setBucketMultiplier(
      /^\d+$/.test(value) ? parseInt(value) : "AUTO",
    );
  }
}
