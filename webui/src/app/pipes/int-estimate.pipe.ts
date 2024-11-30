import { inject, Pipe, PipeTransform } from "@angular/core";
import { TranslocoService } from "@jsverse/transloco";

@Pipe({
  name: "intEstimate",
  standalone: true,
  pure: false,
})
export class IntEstimatePipe implements PipeTransform {
  transloco = inject(TranslocoService);

  transform(n: number, isEstimate = true, sigFigs = 2): string {
    if (isEstimate && n > 0 && sigFigs > 0) {
      const magnitude = Math.floor(Math.log10(Math.abs(n)));
      const scale = Math.pow(10, magnitude - (sigFigs - 1));
      n = Math.round(n / scale) * scale;
    }
    const str = Intl.NumberFormat(this.transloco.getActiveLang()).format(n);
    return isEstimate ? `~${str}` : str;
  }
}
