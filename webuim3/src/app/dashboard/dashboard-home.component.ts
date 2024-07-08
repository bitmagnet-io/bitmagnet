import { Component, inject } from '@angular/core';
import { Breakpoints, BreakpointObserver } from '@angular/cdk/layout';
import { map } from 'rxjs/operators';
import { AsyncPipe } from '@angular/common';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatMenuModule } from '@angular/material/menu';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { HealthModule } from '../health/health.module';
import {ChartComponent} from "../charting/chart.component";
import {queueChartAdapterTotals} from "../queue/queue-chart-adapter.totals";
import {QueueMetricsController} from "../queue/queue-metrics.controller";
import {QueueModule} from "../queue/queue.module";
import {queueChartAdapterTimeline} from "../queue/queue-chart-adapter.timeline";
import {Apollo} from "apollo-angular";
import {TranslocoDirective} from "@jsverse/transloco";
import {MatToolbar} from "@angular/material/toolbar";
import {MatDivider} from "@angular/material/divider";

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard-home.component.html',
  styleUrl: './dashboard-home.component.scss',
  standalone: true,
  imports: [
    AsyncPipe,
    MatGridListModule,
    MatMenuModule,
    MatIconModule,
    MatButtonModule,
    MatCardModule,
    HealthModule,
    ChartComponent,
    QueueModule,
    TranslocoDirective,
    MatToolbar,
    MatDivider,
  ],
})
export class DashboardHomeComponent {
  private apollo = inject(Apollo);
  private breakpointObserver = inject(BreakpointObserver);
  queueMetricsController = new QueueMetricsController(
    this.apollo,
    {
      buckets: {
        duration: "minute",
        multiplier: 5,
        timeframe: "days_1"
      },
      autoRefresh: "off",
    })

  totals = queueChartAdapterTotals
  timeline = queueChartAdapterTimeline

  /** Based on the screen size, switch from standard to one column per row */
  cards = this.breakpointObserver.observe(Breakpoints.Handset).pipe(
    map(({ matches }) => {
      return [
        { title: 'Card 1', cols: 1, rows: 1 },
        { title: 'Card 2', cols: matches ? 1 : 2, rows: 1 },
        { title: 'Card 3', cols: 1, rows: matches ? 1 : 2 },
        { title: 'Card 4', cols: 1, rows: 1 },
      ];
    }),
  );
}
