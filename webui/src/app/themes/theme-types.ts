import {
  themeBaseColors,
  themeColorHues,
  themeColors,
} from "./theme-constants";

export type ThemeBaseColor = (typeof themeBaseColors)[number];

export type ThemeColorHue = (typeof themeColorHues)[number];

export type ThemeColor = (typeof themeColors)[number];

export type ThemeColors = Record<ThemeColor, string>;

export type ThemeType = "light" | "dark";

export type ThemeInfo = {
  colors: ThemeColors;
  type: ThemeType;
};
