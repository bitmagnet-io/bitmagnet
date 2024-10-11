import { ChartConfiguration, ChartType } from "chart.js";
import { ThemeColors } from "../themes/theme-types";

export type ChartDependencies = {
  colors: ThemeColors;
};

export type ChartConfigFactory<
  Data = unknown,
  Type extends ChartType = ChartType,
> = (data: Data | undefined) => ChartConfiguration<Type>;

export interface ChartAdapter<
  Data = unknown,
  Type extends ChartType = ChartType,
> {
  create: ChartConfigFactory<Data, Type>;
}
