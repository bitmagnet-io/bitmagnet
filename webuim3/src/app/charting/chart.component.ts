import {Component, inject, Input, OnInit, ViewChild} from '@angular/core';
import {BaseChartDirective} from "ng2-charts";
import {MatButton} from "@angular/material/button";
import {ChartAdapter} from "./types";
import {ChartConfiguration, ChartType} from "chart.js";
import {Observable} from "rxjs";

@Component({
  selector: 'app-chart',
  standalone: true,
  imports: [BaseChartDirective, MatButton],
  templateUrl: './chart.component.html',
  styleUrl: './chart.component.scss'
})
export class ChartComponent<Data = unknown, Type extends ChartType = ChartType> implements OnInit{
  @Input() $data: Observable<any>
  @Input() adapter: ChartAdapter<Data, Type>;
  @Input() width: number = 500;
  @Input() height: number = 500;

  chartConfig: ChartConfiguration

  ngOnInit() {
    this.chartConfig = this.adapter.create() as ChartConfiguration
    this.$data.subscribe((data) => {
      this.chartConfig = this.adapter.create(data) as ChartConfiguration
    })
  }
}


