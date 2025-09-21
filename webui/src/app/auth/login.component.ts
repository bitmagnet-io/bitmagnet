import { Component, inject } from "@angular/core";
import { FormControl } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
import { EMPTY, map, take } from "rxjs";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { TranslocoService } from "@jsverse/transloco";
import { ErrorsService } from "../errors/errors.service";
import { AppModule } from "../app.module";
import { AuthService } from "./auth.service";

@Component({
  selector: "app-login",
  template: `
    <ng-container *transloco="let t">
      <app-document-title [parts]="[t('auth.login')]" />
      <mat-card>
        <mat-card-header>
          <h2>{{ t("auth.login") }}</h2>
        </mat-card-header>
        <mat-card-content>
          <mat-form-field>
            <mat-label>{{ t("auth.username") }}</mat-label>
            <input
              matInput
              [formControl]="usernameCtrl"
              (keyup.enter)="login()"
            />
          </mat-form-field>
          <mat-form-field>
            <mat-label>{{ t("auth.password") }}</mat-label>
            <input
              matInput
              type="password"
              [formControl]="passwordCtrl"
              (keyup.enter)="login()"
            />
          </mat-form-field>
        </mat-card-content>
        <mat-card-actions class="button-row">
          <button mat-stroked-button (click)="login()">
            <mat-icon>login</mat-icon>{{ t("auth.login") }}
          </button>
          <a mat-stroked-button routerLink="/register">
            <mat-icon>person_add</mat-icon>{{ t("auth.register") }}
          </a>
        </mat-card-actions>
      </mat-card>
    </ng-container>
  `,
  styles: [
    `
      mat-card {
        max-width: 400px;
        margin: 30px auto;
      }
      mat-card-content {
        display: flex;
        flex-direction: column;
        gap: 1em;
      }
      button {
        margin-right: 10px;
      }
    `,
  ],
  standalone: true,
  imports: [AppModule],
})
export class LoginComponent {
  private authService = inject(AuthService);
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private errors = inject(ErrorsService);
  private transloco = inject(TranslocoService);
  usernameCtrl = new FormControl<string>("");
  passwordCtrl = new FormControl<string>("");

  constructor() {
    this.authService.resetLoginError();
    this.authService.loginError$
      .pipe(takeUntilDestroyed())
      .subscribe(
        (error) =>
          error &&
          this.errors.addError(
            `${this.transloco.translate("auth.login_failed")}: ${error}`,
          ),
      );
  }

  login() {
    this.authService
      .login(this.usernameCtrl.value!, this.passwordCtrl.value!)
      .pipe(
        map(() => {
          void this.router.navigate(
            [this.route.snapshot.queryParams["returnUrl"] || "/"],
            { replaceUrl: true },
          );
          return EMPTY;
        }),
        take(1),
      )
      .subscribe();
  }
}
