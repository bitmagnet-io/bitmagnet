import { inject, Injectable } from '@angular/core';
import dayjs, { Dayjs } from 'dayjs';
import duration, { DurationUnitsObjectType } from 'dayjs/plugin/duration.js';
import relativeTime from 'dayjs/plugin/relativeTime.js';
import { TranslateManager } from '../i18n/translate-manager.service';
import 'dayjs/locale/ar.js';
import 'dayjs/locale/de.js';
import 'dayjs/locale/en.js';
import 'dayjs/locale/fr.js';
import 'dayjs/locale/zh.js';

dayjs.extend(duration);
dayjs.extend(relativeTime);

@Injectable({ providedIn: 'root' })
export class DayjsService {
  private translate = inject(TranslateManager);

  createDate(
    date?: dayjs.ConfigType,
    format?: dayjs.OptionType,
    locale?: string,
  ): Dayjs {
    return dayjs(date, format).locale(
      locale ?? this.translate.transloco.getActiveLang(),
    );
  }

  createDuration(d: DurationUnitsObjectType, locale?: string) {
    return dayjs
      .duration(d)
      .locale(locale ?? this.translate.transloco.getActiveLang());
  }
}
