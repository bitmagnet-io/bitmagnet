import { Component, inject, Input, OnInit } from "@angular/core";
import { BaseChartDirective } from "ng2-charts";
import { ChartConfiguration, ChartType } from "chart.js";
import { Observable } from "rxjs";
import { TranslocoDirective, TranslocoService } from "@jsverse/transloco";
import {
  MatCard,
  MatCardContent,
  MatCardHeader,
  MatCardTitle,
} from "@angular/material/card";
import { MatIcon } from "@angular/material/icon";
import { MatTooltip } from "@angular/material/tooltip";
import { ThemeInfoService } from "../themes/theme-info.service";
import { ChartAdapter } from "./types";

@Component({
  selector: "app-chart",
  standalone: true,
  imports: [
    BaseChartDirective,
    MatCard,
    MatCardContent,
    MatCardHeader,
    MatCardTitle,
    TranslocoDirective,
    MatIcon,
    MatTooltip,
  ],
  templateUrl: "./chart.component.html",
  styleUrl: "./chart.component.scss",
})
export class ChartComponent<Data = unknown, Type extends ChartType = ChartType>
  implements OnInit
{
  private themeInfo = inject(ThemeInfoService);
  private transloco = inject(TranslocoService);

  @Input() title: string;
  @Input() $data: Observable<Data> = new Observable();
  @Input() adapter: ChartAdapter<Data, Type>;
  @Input() width: number = 500;
  @Input() height: number = 500;

  chartConfig: ChartConfiguration;

  private data: Data;

  protected legend = true;

  ngOnInit() {
    this.updateChart();
    this.$data.subscribe((data) => {
      this.data = data;
      this.updateChart();
    });
    this.themeInfo.info$.subscribe(() => {
      this.updateChart();
    });
    this.transloco.langChanges$.subscribe(() => {
      this.updateChart();
    });
  }

  protected toggleLegend() {
    this.legend = !this.legend;
    this.updateChart();
  }

  private updateChart() {
    this.chartConfig = this.adapter.create(this.data, {
      legend: this.legend,
    }) as ChartConfiguration;
  }
}
