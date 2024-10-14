import { ThemeBaseColor, ThemeColor, ThemeColorHue } from "./theme-types";

export const createThemeColor = (
  base: ThemeBaseColor,
  hue: ThemeColorHue,
): ThemeColor => `${base}-${hue}`;
