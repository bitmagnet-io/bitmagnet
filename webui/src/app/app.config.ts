import {
  ApplicationConfig,
  provideZoneChangeDetection,
  inject,
} from "@angular/core";
import { provideRouter, withComponentInputBinding } from "@angular/router";

import { provideAnimationsAsync } from "@angular/platform-browser/animations/async";
import {
  HttpHeaders,
  provideHttpClient,
  withInterceptorsFromDi,
} from "@angular/common/http";
import { setContext } from "@apollo/client/link/context";
import { provideTransloco, TranslocoService } from "@jsverse/transloco";
import { provideCharts, withDefaultRegisterables } from "ng2-charts";
import { provideApollo } from "apollo-angular";
import { HttpBatchLink } from "apollo-angular/http";
import { ApolloLink, InMemoryCache } from "@apollo/client/core";
import { graphqlEndpoint } from "../environments/environment";
import { TranslocoImportLoader } from "./i18n/transloco.loader";
import { routes } from "./app.routes";
import { AuthTokenService } from "./auth/auth-token.service";

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes, withComponentInputBinding()),
    provideAnimationsAsync("animations"),
    provideHttpClient(withInterceptorsFromDi()),
    HttpBatchLink,
    provideApollo(() => {
      const httpLink = inject(HttpBatchLink);
      const transloco = inject(TranslocoService);
      const tokenService = inject(AuthTokenService);

      const middleware = new ApolloLink((operation, forward) => {
        let headers = new HttpHeaders().set(
          "Accept-Language",
          transloco.getActiveLang(),
        );

        const token = tokenService.getToken();

        if (token) {
          headers = headers.set("Authorization", "Bearer " + token);
        }

        operation.setContext({
          headers,
        });
        return forward(operation);
      });

      return {
        link: middleware.concat(
          httpLink.create({
            uri: graphqlEndpoint,
            batchMax: 10,
            batchInterval: 50,
          }),
        ),
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
