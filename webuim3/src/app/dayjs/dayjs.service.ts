import {inject, Injectable} from "@angular/core";
import dayjs, {Dayjs} from "dayjs";
import duration, {CreateDurationType, DurationUnitsObjectType} from "dayjs/plugin/duration"
import relativeTime from  "dayjs/plugin/relativeTime"
import {TranslateManager} from "../i18n/translate-manager.service";
import "dayjs/locale/ar"
import "dayjs/locale/de"
import "dayjs/locale/en"
import "dayjs/locale/fr"
import "dayjs/locale/zh"

dayjs.extend(duration)
dayjs.extend(relativeTime)

@Injectable({providedIn: "root"})
export class DayjsService {
  private translate = inject(TranslateManager)

  createDate(date?: dayjs.ConfigType, format?: dayjs.OptionType, locale?: string): Dayjs {
    return dayjs(date, format).locale(locale ?? this.translate.transloco.getActiveLang())
  }

  createDuration(d: DurationUnitsObjectType, locale?: string) {
    return dayjs.duration(d).locale(locale ?? this.translate.transloco.getActiveLang())
  }
}
