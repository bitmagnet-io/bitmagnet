import { Injectable } from '@angular/core';
import { Translation, TranslocoLoader } from '@jsverse/transloco';
import i18n from "./translations";

@Injectable({ providedIn: 'root' })
export class TranslocoImportLoader implements TranslocoLoader {
  getTranslation(lang: string) {
    if (lang in i18n) {
      return Promise.resolve(i18n[lang as keyof typeof i18n] as Translation).then(stripMissing);
    } else {
      return Promise.reject(`Translation not found: ${lang}`);
    }
  }
}

const missingValues = ["__missing__", "__fallback__"]

const stripMissing = (tr: Translation) =>
  Object.fromEntries(Object.entries(tr).flatMap(([k, v]) => {
    if (typeof v === "object") {
      v = stripMissing(v);
    } else if (typeof v === "string" && missingValues.includes(v)) {
      return [];
    }
    return [[k, v]]
  }));
