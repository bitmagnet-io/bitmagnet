import {
  ChartConfiguration,
  ChartType,
  ChartEvent,
  LegendItem,
  LegendElement,
} from "chart.js";

export type LegendOnClick = (
  e: ChartEvent,
  item: LegendItem,
  ctx: LegendElement<"line">,
) => void;

export type FactoryParams = {
  legend: boolean;
  hiddenDatasets: Map<string, boolean>;
  legendOnClick: LegendOnClick;
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
