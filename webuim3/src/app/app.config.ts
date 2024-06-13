import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import {provideRouter, withComponentInputBinding} from '@angular/router';

import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import {
  provideHttpClient,
  withInterceptorsFromDi,
} from '@angular/common/http';
import { routes } from './app.routes';
import {provideAnimations} from "@angular/platform-browser/animations";

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes, withComponentInputBinding()),
    provideAnimationsAsync('animations'),
    // provideAnimations(),
    provideHttpClient(withInterceptorsFromDi()),
  ],
};
