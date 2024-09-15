import {inject, Pipe, PipeTransform} from '@angular/core';
import {formatTimeAgo} from "./dates.utils";
import {DateType} from "./dates.types";
import {TranslocoService} from "@jsverse/transloco";

@Pipe({
  name: 'timeAgo',
  standalone: true,
})
export class TimeAgoPipe implements PipeTransform {
  transloco = inject(TranslocoService);

  transform(value: DateType): string {
    return formatTimeAgo(value, this.transloco.getActiveLang())
  }
}
