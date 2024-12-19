import { Component, inject, Input, OnInit, ViewChild } from "@angular/core";
import { BaseChartDirective } from "ng2-charts";
import {
  ChartConfiguration,
  ChartType,
  LegendItem,
  ChartEvent,
} from "chart.js";
import { Observable } from "rxjs";
import { TranslocoService } from "@jsverse/transloco";
import { ThemeInfoService } from "../themes/theme-info.service";
import { BreakpointsService } from "../layout/breakpoints.service";
import { AppModule } from "../app.module";
import { ChartAdapter } from "./types";

@Component({
  selector: "app-chart",
  standalone: true,
  imports: [AppModule, BaseChartDirective],
  templateUrl: "./chart.component.html",
  styleUrl: "./chart.component.scss",
})
export class ChartComponent<Data = unknown, Type extends ChartType = ChartType>
  implements OnInit
{
  private themeInfo = inject(ThemeInfoService);
  private transloco = inject(TranslocoService);
  breakpoints = inject(BreakpointsService);
  private hiddenDatasets = new Map<string, boolean>();

  @Input() $data: Observable<Data> = new Observable();
  @Input() adapter: ChartAdapter<Data, Type>;
  @Input() width: number = 500;
  @Input() height: number = 500;
  @Input() title: string = "chart title";
  @ViewChild(BaseChartDirective) protected basechart: BaseChartDirective;

  chartConfig: ChartConfiguration;
  chartShowLegend: boolean = true;

  private data: Data;

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
    this.chartShowLegend = !this.chartShowLegend;
    this.legendConfig();
  }

  private setDatasetVisible() {
    for (let i = 0; i < this.basechart.chart!.data.datasets.length; i++) {
      const meta = this.basechart.chart!.getDatasetMeta(i);
      if (this.hiddenDatasets.has(meta.label)) {
        meta.hidden = this.hiddenDatasets.get(meta.label)!;
      }
    }
  }

  private legendOnClick(e: ChartEvent, legendItem: LegendItem) {
    const meta = this.basechart.chart!.getDatasetMeta(legendItem.datasetIndex!);
    this.hiddenDatasets.set(meta.label, !meta.hidden);
    this.setDatasetVisible();
    this.basechart.chart!.update();
  }

  private legendConfig() {
    // there is a delay between chart being accessible and an update
    // being possible hence delay here
    setTimeout(() => {
      this.basechart.chart!.legend!.options.display = this.chartShowLegend;
      this.basechart.chart!.options.plugins!.legend!.onClick =
        this.legendOnClick.bind(this);
      this.setDatasetVisible();
      this.basechart.chart!.update();
    }, 5);
  }

  private updateChart() {
    this.chartConfig = this.adapter.create(this.data) as ChartConfiguration;
    this.legendConfig();
  }
}
