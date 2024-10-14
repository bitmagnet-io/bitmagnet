import { inject, Pipe, PipeTransform } from "@angular/core";
import { TranslocoService } from "@jsverse/transloco";
import { formatTimeAgo } from "../dates/dates.utils";
import { DateType } from "../dates/dates.types";

@Pipe({
  name: "timeAgo",
  standalone: true,
  pure: false,
})
export class TimeAgoPipe implements PipeTransform {
  transloco = inject(TranslocoService);

  transform(value: DateType): string {
    return formatTimeAgo(value, this.transloco.getActiveLang());
  }
}
