import { ar } from "date-fns/locale/ar";
import { de } from "date-fns/locale/de";
import { enUS as en } from "date-fns/locale/en-US";
import { es } from "date-fns/locale/es";
import { fr } from "date-fns/locale/fr";
import { hi } from "date-fns/locale/hi";
import { ja } from "date-fns/locale/ja";
import { nl } from "date-fns/locale/nl";
import { pt } from "date-fns/locale/pt";
import { ru } from "date-fns/locale/ru";
import { tr } from "date-fns/locale/tr";
import { uk } from "date-fns/locale/uk";
import { zhCN as zh } from "date-fns/locale/zh-CN";
import type { Locale } from "date-fns";

const locales: Record<string, Locale> = {
  ar,
  de,
  en,
  es,
  fr,
  hi,
  ja,
  nl,
  pt,
  ru,
  tr,
  uk,
  zh,
};

export default locales;

export const resolveDateLocale = (locale: string): Locale =>
  locales[locale] ?? en;
