import {
  ApplicationConfig,
  provideZoneChangeDetection,
  inject,
} from "@angular/core";
import { provideRouter, withComponentInputBinding } from "@angular/router";

import { provideAnimationsAsync } from "@angular/platform-browser/animations/async";
import {
  provideHttpClient,
  withInterceptorsFromDi,
} from "@angular/common/http";
import { provideTransloco } from "@jsverse/transloco";
import { provideCharts, withDefaultRegisterables } from "ng2-charts";
import { provideApollo } from "apollo-angular";
import { HttpLink } from "apollo-angular/http";
import { InMemoryCache } from "@apollo/client/core";
import { graphqlEndpoint } from "../environments/environment";
import { TranslocoImportLoader } from "./i18n/transloco.loader";
import { routes } from "./app.routes";

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes, withComponentInputBinding()),
    provideAnimationsAsync("animations"),
    provideHttpClient(withInterceptorsFromDi()),
    provideHttpClient(),
    provideApollo(() => {
      const httpLink = inject(HttpLink);
      return {
        link: httpLink.create({ uri: graphqlEndpoint }),
        cache: new InMemoryCache({
          typePolicies: {
            Query: {
              fields: {
                search: {
                  merge(
                    existing: Record<string, unknown>,
                    incoming: Record<string, unknown>,
                  ): Record<string, unknown> {
                    return { ...existing, ...incoming };
                  },
                },
              },
            },
          },
        }),
      };
    }),
    provideTransloco({
      config: {
        availableLangs: [
          {
            id: "ar",
            label: "العربية",
          },
          {
            id: "ca",
            label: "Català",
          },
          {
            id: "de",
            label: "Deutsch",
          },
          {
            id: "en",
            label: "English",
          },
          {
            id: "es",
            label: "Español",
          },
          {
            id: "fr",
            label: "Français",
          },
          {
            id: "hi",
            label: "हिन्दी",
          },
          {
            id: "ja",
            label: "日本語",
          },
          {
            id: "nl",
            label: "Nederlands",
          },
          {
            id: "pt",
            label: "Português",
          },
          {
            id: "ru",
            label: "Русский",
          },
          {
            id: "tr",
            label: "Türkçe",
          },
          {
            id: "uk",
            label: "Українська",
          },
          {
            id: "zh",
            label: "中文",
          },
        ],
        defaultLang: "en",
        fallbackLang: "en",
        missingHandler: {
          // It will use the first language set in the `fallbackLang` property
          useFallbackTranslation: true,
        },
        // Remove this option if your application doesn't support changing language in runtime.
        reRenderOnLangChange: true,
        prodMode: false,
      },
      loader: TranslocoImportLoader,
    }),
    provideCharts(withDefaultRegisterables()),
  ],
};
