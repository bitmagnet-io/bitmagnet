import {formatDistance, formatDistanceToNow} from "date-fns";
import {resolveDateLocale} from "./dates.locales";
import {DateType} from "./dates.types";

export const formatTimeAgo = (date: DateType, locale: string) =>
  formatDistanceToNow(date, {
    addSuffix: true,
    locale: resolveDateLocale(locale),
  })

export const formatDuration = (seconds: number, locale: string) =>
  formatDistance(seconds * 1000, 0, {
    includeSeconds: true,
    locale: resolveDateLocale(locale),
  })
