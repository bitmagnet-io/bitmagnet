import { inject, Pipe, PipeTransform } from '@angular/core';
import { DayjsService } from '../dayjs/dayjs.service';

@Pipe({
  name: 'humanTime',
  standalone: true,
})
export class HumanTimePipe implements PipeTransform {
  private dayjs = inject(DayjsService);

  transform(value: Date | string, locale?: string): string {
    return this.dayjs
      .createDate(undefined, undefined, locale)
      .to(this.dayjs.createDate(value, undefined, locale));
  }
}
