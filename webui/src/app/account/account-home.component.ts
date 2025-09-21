import { Component, inject } from "@angular/core";
import { HealthModule } from "../health/health.module";
import { AppModule } from "../app.module";
import { DocumentTitleComponent } from "../layout/document-title.component";
import { BreakpointsService } from "../layout/breakpoints.service";
import { map } from "rxjs";
import { AuthService } from "../auth/auth.service";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";

@Component({
  selector: "app-account-home",
  template: `
    <ng-container *transloco="let t">
      <app-document-title [parts]="[t('routes.account')]" />
      @if (user$ | async; as user) {
        <mat-card class="account-card">
          <mat-card-header>
            <mat-toolbar>
              <h2>
                <mat-icon>person</mat-icon>{{ t("routes.account") }}:
                {{ user.username }}
              </h2>
            </mat-toolbar>
          </mat-card-header>
          <mat-card-content> </mat-card-content>
        </mat-card>
      }
    </ng-container>
  `,
  styles: [
    `
      .grid-container {
        margin: 20px;
      }

      .more-button {
        position: absolute;
        top: 5px;
        right: 10px;
      }

      app-health-card {
        width: 100%;
        height: 100%;
        mat-card {
          height: 100%;
        }
      }

      mat-grid-tile mat-card {
        width: 100%;
      }

      mat-toolbar h2 mat-icon {
        position: relative;
        top: 3px;
        margin-right: 14px;
        margin-left: 32px;
      }
    `,
  ],
  standalone: true,
  imports: [AppModule, HealthModule, DocumentTitleComponent],
})
export class AccountHomeComponent {
  breakpoints = inject(BreakpointsService);
  private authService = inject(AuthService);

  user$ = this.authService.self$.pipe(
    takeUntilDestroyed(),
    map((self) => self.user),
  );
}
