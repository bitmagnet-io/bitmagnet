import { ThemeColors, ThemeInfo } from "./theme-types";

export const themeBaseColors = [
  "primary",
  "secondary",
  "tertiary",
  "neutral",
  "neutral-variant",
  "error",
  "caution",
  "success",
] as const;

export const themeColorHues = [20, 40, 50, 60, 80] as const;

export const themeColors = [
  "background",
  "foreground",
  ...themeBaseColors.flatMap((baseColor) =>
    themeColorHues.map(
      (
        n,
      ): `${(typeof themeBaseColors)[number]}-${(typeof themeColorHues)[number]}` =>
        `${baseColor}-${n}`,
    ),
  ),
] as const;

export const emptyThemeInfo: ThemeInfo = {
  type: "light",
  colors: {} as ThemeColors,
};
