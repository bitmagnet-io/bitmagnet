import {Component, inject, ViewChild} from '@angular/core';
import {BaseChartDirective} from "ng2-charts";
import {ChartConfiguration, ChartData, ChartEvent} from "chart.js";
import {MatButton} from "@angular/material/button";
import {QueueMetricsService} from "./queue-metrics.service";
import {map} from "rxjs/operators";

@Component({
  selector: 'app-queue-chart',
  standalone: true,
  imports: [BaseChartDirective, MatButton],
  providers: [QueueMetricsService],
  templateUrl: './queue-chart.component.html',
  styleUrl: './queue-chart.component.scss'
})
export class QueueChartComponent {
  private metricsService = inject(QueueMetricsService)

  @ViewChild(BaseChartDirective) chart?: BaseChartDirective;

  public barChartOptions: ChartConfiguration<'bar'>['options'] = {
    // We use these empty structures as placeholders for dynamic theming.
    scales: {
      x: {},
      y: {
        // min: 10,
      },
    },
    plugins: {
      legend: {
        display: true,
      },
      // datalabels: {
      //   anchor: 'end',
      //   align: 'end',
      // },
    },
  };
  public barChartType = 'bar' as const;

  public barChartData: ChartData<'bar'> = {
    labels: ['2006', '2007', '2008', '2009', '2010', '2011', '2012'],
    datasets: [
      { data: [65, 59, 80, 81, 56, 55, 40], label: 'Series A' },
      { data: [28, 48, 40, 19, 86, 27, 90], label: 'Series B' },
    ],
  };

  // events
  public chartClicked({
                        event,
                        active,
                      }: {
    event?: ChartEvent;
    active?: object[];
  }): void {
    console.log(event, active);
    this.metricsService.request({bucketDuration: "minute"}).pipe(map(console.log)).subscribe()
  }

  public chartHovered({
                        event,
                        active,
                      }: {
    event?: ChartEvent;
    active?: object[];
  }): void {
    // console.log(event, active);
  }

  public randomize(): void {
    // Only Change 3 values
    this.barChartData.datasets[0].data = [
      Math.round(Math.random() * 100),
      59,
      80,
      Math.round(Math.random() * 100),
      56,
      Math.round(Math.random() * 100),
      40,
    ];

    this.chart?.update();
  }
}
