import { Pipe, PipeTransform } from '@angular/core';
import {formatTimeAgo} from "./dates.utils";
import {DateType} from "./dates.types";

@Pipe({
  name: 'timeAgo',
  standalone: true,
})
export class TimeAgoPipe implements PipeTransform {
  transform(value: DateType, locale: string): string {
    return formatTimeAgo(value, locale)
  }
}
