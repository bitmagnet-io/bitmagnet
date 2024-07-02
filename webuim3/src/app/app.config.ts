import {
  ApplicationConfig,
  provideZoneChangeDetection,
} from '@angular/core';
import { provideRouter, withComponentInputBinding } from '@angular/router';

import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import {
  provideHttpClient,
  withInterceptorsFromDi,
} from '@angular/common/http';
import { provideTransloco } from '@jsverse/transloco';
import { routes } from './app.routes';
import { TranslocoImportLoader } from './i18n/transloco.loader';
import { provideCharts, withDefaultRegisterables } from 'ng2-charts';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes, withComponentInputBinding()),
    provideAnimationsAsync('animations'),
    provideHttpClient(withInterceptorsFromDi()),
    provideHttpClient(),
    provideTransloco({
      config: {
        availableLangs: [
          {
            id: 'ar',
            label: 'العربية',
          },
          {
            id: 'en',
            label: 'English',
          },
          {
            id: 'fr',
            label: 'Français',
          },
          {
            id: 'es',
            label: 'Español',
          },
          {
            id: 'de',
            label: 'Deutsch',
          },
          {
            id: 'zh',
            label: '中文',
          },
        ],
        defaultLang: 'en',
        fallbackLang: 'en',
        missingHandler: {
          // It will use the first language set in the `fallbackLang` property
          useFallbackTranslation: true,
        },
        // Remove this option if your application doesn't support changing language in runtime.
        reRenderOnLangChange: true,
        prodMode: false,
      },
      loader: TranslocoImportLoader,
    }), provideCharts(withDefaultRegisterables()),
  ],
};
