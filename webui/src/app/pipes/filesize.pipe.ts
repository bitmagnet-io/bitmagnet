import { Pipe, PipeTransform } from '@angular/core';
import { filesize } from 'filesize';
@Pipe({
  name: 'filesize',
  standalone: true,
})
export class FilesizePipe implements PipeTransform {
  transform(value: number, locale?: string): string {
    return filesize(value, { locale });
  }
}
