import { inject, Pipe, PipeTransform } from "@angular/core";
import { filesize } from "filesize";
import { TranslocoService } from "@jsverse/transloco";

@Pipe({
  name: "filesize",
  standalone: true,
  pure: false,
})
export class FilesizePipe implements PipeTransform {
  transloco = inject(TranslocoService);

  transform(value: number, base: 2 | 10 = 2): string {
    return filesize(value, { locale: this.transloco.getActiveLang(), base });
  }
}
