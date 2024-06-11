import { DOCUMENT } from '@angular/common';
import { Injectable, inject } from '@angular/core';
// import { BehaviorSubject } from 'rxjs/internal/BehaviorSubject';
// import { take } from 'rxjs/operators';
import { BrowserStorageService } from '../browser-storage/browser-storage.service';
import {defaultDarkTheme, defaultLightTheme, ThemeKey, themes} from "./theme-registry";

const LOCAL_STORAGE_KEY = 'bitmagnet-theme';

@Injectable({ providedIn: 'root' })
export class ThemeManager {
  private document = inject(DOCUMENT);
  private browserStorage = inject(BrowserStorageService);
  // private _isDarkSub = new BehaviorSubject(false);
  // isDark$ = this._isDarkSub.asObservable();
  private _window = this.document.defaultView;
  public selectedTheme?: ThemeKey;
  public themes = Object.values(themes)

  constructor() {
    this.setTheme(this.getPreferredTheme());
      this.windowMatchMediaPrefersDark()?.addEventListener('change', () => {
          const storedTheme = this.getStoredTheme();
          if (!storedTheme) {
            this.setTheme(this.getAutoTheme());
          }
        });
  }

  getStoredTheme = (): ThemeKey | undefined => {
    const value = this.browserStorage.get(LOCAL_STORAGE_KEY)
    return value && value in themes ? value as ThemeKey : undefined
  }

  setStoredTheme = (theme: ThemeKey | "auto") => {
    if (theme === "auto") {
      this.browserStorage.remove(LOCAL_STORAGE_KEY)
    } else {
      this.browserStorage.set(LOCAL_STORAGE_KEY, theme);
    }
  };

  getPreferredTheme = (): ThemeKey => {
    return this.getStoredTheme() ?? this.getAutoTheme();
  };

  getAutoTheme = (): ThemeKey => {
    return this.windowMatchMediaPrefersDark()?.matches
      ? defaultDarkTheme
      : defaultLightTheme;
  };

  setTheme = (theme: string) => {
    if (theme === "auto" || !(theme in themes)) {
      theme = this.getAutoTheme()
      this.selectedTheme = undefined;
    } else {
      this.selectedTheme = theme as ThemeKey;
    }
    this.document.documentElement.setAttribute('data-bitmagnet-theme', theme);
    this.setStoredTheme(this.selectedTheme ?? 'auto');
  };

  // setMaterialTheme() {
  //   this.isDark$.pipe(take(1)).subscribe((isDark) => {
  //     if (isDark) {
  //       const href = 'dark-theme.css';
  //       getLinkElementForKey('dark-theme').setAttribute('href', href);
  //       this.document.documentElement.classList.add('dark-theme');
  //     } else {
  //       this.removeStyle('dark-theme');
  //       this.document.documentElement.classList.remove('dark-theme');
  //     }
  //   });
  // }

  removeStyle(key: string) {
    const existingLinkElement = getExistingLinkElementByKey(key);
    if (existingLinkElement) {
      this.document.head.removeChild(existingLinkElement);
    }
  }

  changeTheme(theme: ThemeKey | "auto") {
    this.setStoredTheme(theme);
    this.setTheme(theme);
  }

  windowMatchMediaPrefersDark(): MediaQueryList | undefined {
    return this._window && this._window.matchMedia ? this._window.matchMedia('(prefers-color-scheme: dark)') : undefined
  }
}

function getLinkElementForKey(key: string) {
  return getExistingLinkElementByKey(key) || createLinkElementWithKey(key);
}

function getExistingLinkElementByKey(key: string) {
  return document.head.querySelector(
    `link[rel="stylesheet"].${getClassNameForKey(key)}`
  );
}

function createLinkElementWithKey(key: string) {
  const linkEl = document.createElement('link');
  linkEl.setAttribute('rel', 'stylesheet');
  linkEl.classList.add(getClassNameForKey(key));
  document.head.appendChild(linkEl);
  return linkEl;
}

function getClassNameForKey(key: string) {
  return `style-manager-${key}`;
}
