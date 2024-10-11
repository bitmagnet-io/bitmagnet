export type ThemeInfo<Key extends string = string> = {
  key: Key;
  label: string;
  dark: boolean;
};

export type Themes<Keys extends string = string> = {
  [key in Keys]: ThemeInfo<key>;
};

const _themes = {
  classic: {
    key: "classic" as const,
    label: "Classic",
    dark: false,
  },
  tundra: {
    key: "tundra" as const,
    label: "Tundra",
    dark: true,
  },
};

export type ThemeKey = keyof typeof _themes;

export const themes: Themes<ThemeKey> = _themes;

export const defaultLightTheme = "classic";
export const defaultDarkTheme = "tundra";
