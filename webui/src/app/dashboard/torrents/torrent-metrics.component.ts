import { Component, inject, OnDestroy } from "@angular/core";
import { Apollo } from "apollo-angular";
import { GraphQLModule } from "../../graphql/graphql.module";
import { ChartComponent } from "../../charting/chart.component";
import { BreakpointsService } from "../../layout/breakpoints.service";
import { ErrorsService } from "../../errors/errors.service";
import { AppModule } from "../../app.module";
import {
  autoRefreshIntervalNames,
  defaultBucketParams,
  eventNames,
  resolutionNames,
  timeframeNames,
} from "./torrent-metrics.constants";
import { TorrentMetricsController } from "./torrent-metrics.controller";
import { TorrentChartAdapterTimeline } from "./torrent-chart-adapter.timeline";

@Component({
  selector: "app-torrent-metrics",
  standalone: true,
  templateUrl: "./torrent-metrics.component.html",
  styleUrl: "./torrent-metrics.component.scss",
  imports: [AppModule, ChartComponent, GraphQLModule],
})
export class TorrentMetricsComponent implements OnDestroy {
  breakpoints = inject(BreakpointsService);
  private apollo = inject(Apollo);
  torrentMetricsController = new TorrentMetricsController(
    this.apollo,
    {
      buckets: defaultBucketParams,
      autoRefresh: "seconds_30",
    },
    inject(ErrorsService),
  );
  protected readonly timeline = inject(TorrentChartAdapterTimeline);

  protected readonly resolutionNames = resolutionNames;
  protected readonly timeframeNames = timeframeNames;
  protected readonly autoRefreshIntervalNames = autoRefreshIntervalNames;

  ngOnDestroy() {
    this.torrentMetricsController.setAutoRefreshInterval("off");
  }

  protected readonly eventNames = eventNames;

  handleMultiplierEvent(event: Event) {
    const value = (event.currentTarget as HTMLInputElement).value;
    this.torrentMetricsController.setBucketMultiplier(
      /^\d+$/.test(value) ? parseInt(value) : "AUTO",
    );
  }
}
