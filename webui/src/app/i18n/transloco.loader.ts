import { Injectable } from "@angular/core";
import { Translation, TranslocoLoader } from "@jsverse/transloco";
import i18n from "./translations";

@Injectable({ providedIn: "root" })
export class TranslocoImportLoader implements TranslocoLoader {
  async getTranslation(lang: string) {
    if (lang in i18n) {
      const tr = i18n[lang as keyof typeof i18n] as Translation;
      return stripMissing(tr);
    } else {
      return Promise.reject(new Error(`Translation not found: ${lang}`));
    }
  }
}

const missingValues = ["__missing__", "__fallback__"];

const stripMissing = (tr: Translation): Translation =>
  Object.fromEntries(
    Object.entries(tr).flatMap(([k, v]) => {
      if (typeof v === "object") {
        v = stripMissing(v as Translation);
      } else if (typeof v === "string" && missingValues.includes(v)) {
        return [];
      }
      return [[k, v]];
    }),
  );
