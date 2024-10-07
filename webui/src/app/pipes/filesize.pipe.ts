import { inject, Pipe, PipeTransform } from '@angular/core';
import { filesize } from 'filesize';
import { TranslocoService } from '@jsverse/transloco';
@Pipe({
  name: 'filesize',
  standalone: true,
})
export class FilesizePipe implements PipeTransform {
  transloco = inject(TranslocoService);

  transform(value: number): string {
    return filesize(value, { locale: this.transloco.getActiveLang(), base: 2 });
  }
}
