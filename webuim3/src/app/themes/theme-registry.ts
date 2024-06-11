export type ThemeInfo<Key extends string = string> = {
  key: Key;
  label: string;
  dark: boolean;
};

export type Themes<Keys extends string = string> = {
  [key: string]: ThemeInfo<typeof key>;
};

const _themes = {
  classic: {
    key: "classic",
    label: "Classic",
    dark: false,
  },
  tundra: {
    key: "tundra",
    label: "Tundra",
    dark: true,
  }
}

export type ThemeKey = keyof typeof _themes;

export const themes: Themes<ThemeKey> = _themes;

export const defaultLightTheme = "classic" as const;
export const defaultDarkTheme = "tundra" as const;
