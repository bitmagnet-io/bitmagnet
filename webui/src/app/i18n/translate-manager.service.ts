import { inject, Injectable } from "@angular/core";
import { LangDefinition, TranslocoService } from "@jsverse/transloco";
import { BrowserStorageService } from "../browser-storage/browser-storage.service";

const LOCAL_STORAGE_KEY = "bitmagnet-language";

@Injectable({ providedIn: "root" })
export class TranslateManager {
  transloco = inject(TranslocoService);
  private browserStorage = inject(BrowserStorageService);

  availableLanguages = this.transloco.getAvailableLangs() as LangDefinition[];

  constructor() {
    this.transloco.setActiveLang(this.getPreferredLanguage());
  }

  getPreferredLanguage(): string {
    return this.getStoredLanguage() ?? this.getAutoLanguage();
  }

  getStoredLanguage(): string | undefined {
    const value = this.browserStorage.get(LOCAL_STORAGE_KEY);
    return value && this.transloco.isLang(value) ? value : undefined;
  }

  getAutoLanguage(): string {
    const navLang = navigator?.language?.split("-")?.[0];
    return this.transloco.isLang(navLang) ? navLang : "en";
  }

  setLanguage(lang: string) {
    this.transloco.setActiveLang(lang);
    this.browserStorage.set(LOCAL_STORAGE_KEY, lang);
  }
}
