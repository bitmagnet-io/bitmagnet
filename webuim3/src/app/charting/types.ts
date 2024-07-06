import {ChartConfiguration,ChartType} from "chart.js";

export type ChartConfigFactory<Data = unknown, Type extends ChartType = ChartType> = (data?: Data) => ChartConfiguration<Type>;

export interface ChartAdapter<Data = unknown, Type extends ChartType = ChartType> {
  create: ChartConfigFactory<Data, Type>
}
