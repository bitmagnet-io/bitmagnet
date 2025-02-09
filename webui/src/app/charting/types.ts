import { ChartConfiguration, ChartType } from "chart.js";

export type FactoryParams = {
  legend: boolean;
};

export type ChartConfigFactory<
  Data = unknown,
  Type extends ChartType = ChartType,
> = (data: Data | undefined, params: FactoryParams) => ChartConfiguration<Type>;

export interface ChartAdapter<
  Data = unknown,
  Type extends ChartType = ChartType,
> {
  create: ChartConfigFactory<Data, Type>;
}
