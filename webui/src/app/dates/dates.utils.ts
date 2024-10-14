import {
  formatDistanceToNow,
  formatDuration as fnsFormatDuration,
  Duration,
} from "date-fns";
import { resolveDateLocale } from "./dates.locales";
import { DateType } from "./dates.types";

export const formatTimeAgo = (date: DateType, locale: string) =>
  formatDistanceToNow(date, {
    addSuffix: true,
    locale: resolveDateLocale(locale),
  });

export const formatDuration = (duration: Duration, locale: string) =>
  fnsFormatDuration(duration, {
    locale: resolveDateLocale(locale),
  });
