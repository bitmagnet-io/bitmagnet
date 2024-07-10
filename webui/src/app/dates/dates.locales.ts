import { ar } from 'date-fns/locale/ar'
import { de } from 'date-fns/locale/de'
import { enUS as en } from 'date-fns/locale/en-US'
import { es } from 'date-fns/locale/es'
import { fr } from 'date-fns/locale/fr'
import { zhCN as zh } from 'date-fns/locale/zh-CN'
import type {Locale} from "date-fns";

const locales: Record<string, Locale> = {
  ar,
  de,
  en,
  es,
  fr,
  zh,
}

export default locales

export const resolveDateLocale = (locale: string): Locale =>
  locales[locale as keyof typeof locales] ?? en
