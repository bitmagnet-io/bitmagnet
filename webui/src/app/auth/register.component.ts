import { Component, inject } from "@angular/core";
import { FormControl } from "@angular/forms";
import { Router } from "@angular/router";
import { switchMap, take, Observable, debounceTime } from "rxjs";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { TranslocoService } from "@jsverse/transloco";
import { ErrorStateMatcher } from "@angular/material/core";
import { ErrorsService } from "../errors/errors.service";
import { AppModule } from "../app.module";
import * as generated from "../graphql/generated";
import { AuthService } from "./auth.service";

@Component({
  selector: "app-register",
  template: `
    <ng-container *transloco="let t">
      <app-document-title [parts]="[t('auth.register')]" />
      <mat-card>
        <mat-card-header>
          <h2>{{ t("auth.register") }}</h2>
        </mat-card-header>
        <mat-card-content>
          <mat-form-field>
            <mat-label>{{ t("auth.invitation_code") }}</mat-label>
            <input
              matInput
              [formControl]="invitationCodeCtrl"
              (keyup.enter)="register()"
            />
          </mat-form-field>
          <mat-form-field>
            <mat-label>{{ t("auth.username") }}</mat-label>
            <input
              matInput
              [formControl]="usernameCtrl"
              (keyup.enter)="register()"
            />
          </mat-form-field>
          <mat-form-field>
            <mat-label>{{ t("general.email") }}</mat-label>
            <input
              matInput
              type="email"
              [formControl]="emailCtrl"
              (keyup.enter)="register()"
            />
          </mat-form-field>
          <mat-form-field>
            <mat-label>{{ t("auth.password") }}</mat-label>
            <input
              matInput
              type="password"
              [formControl]="passwordCtrl"
              [errorStateMatcher]="passwordEntropyMatcher"
              (keyup.enter)="register()"
            />
            @if (passwordEntropy$ | async; as passwordEntropy) {
              <mat-progress-bar
                mode="determinate"
                [value]="passwordEntropy.entropy"
              ></mat-progress-bar>
            }
          </mat-form-field>
          @if (
            passwordCtrl.value && !((passwordEntropy$ | async)?.valid ?? true)
          ) {
            <mat-error class="error-entropy">{{
              t("auth.password_has_insufficient_entropy")
            }}</mat-error>
          }
          <mat-form-field>
            <mat-label>{{ t("auth.confirm_password") }}</mat-label>
            <input
              matInput
              type="password"
              [formControl]="confirmPasswordCtrl"
              [errorStateMatcher]="confirmPasswordMatcher"
              (keyup.enter)="register()"
            />
            @if (passwordCtrl.value !== confirmPasswordCtrl.value) {
              <mat-error>{{ t("auth.passwords_do_not_match") }}</mat-error>
            }
          </mat-form-field>
        </mat-card-content>
        <mat-card-actions class="button-row">
          <button mat-stroked-button (click)="register()">
            <mat-icon>person_add</mat-icon>{{ t("auth.register") }}
          </button>
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
      mat-card-actions {
        margin-top: 10px;
      }
      ::ng-deep mat-error {
        font-size: 0.75rem;
      }
      .error-entropy {
        margin-top: -35px;
        margin-left: 16px;
      }
    `,
  ],
  standalone: true,
  imports: [AppModule],
})
export class RegisterComponent {
  private authService = inject(AuthService);
  private router = inject(Router);
  private errors = inject(ErrorsService);
  private transloco = inject(TranslocoService);
  invitationCodeCtrl = new FormControl<string>("");
  usernameCtrl = new FormControl<string>("");
  emailCtrl = new FormControl<string>("");
  passwordCtrl = new FormControl<string>("");
  confirmPasswordCtrl = new FormControl<string>("");

  passwordEntropy$ = this.passwordCtrl.valueChanges.pipe(
    takeUntilDestroyed(),
    debounceTime(100),
    switchMap((password) => this.authService.passwordEntropy(password ?? "")),
  );

  passwordEntropyMatcher = new PasswordEntropyMatcher(this.passwordEntropy$);
  confirmPasswordMatcher = new ConfirmPasswordMatcher(this.passwordCtrl);

  constructor() {
    this.authService.resetRegisterError();
    this.authService.registerError$
      .pipe(takeUntilDestroyed())
      .subscribe(
        (error) =>
          error &&
          this.errors.addError(
            `${this.transloco.translate("auth.register_failed")}: ${error}`,
          ),
      );
  }

  register() {
    if (this.passwordCtrl.value !== this.confirmPasswordCtrl.value) {
      this.errors.addError(
        this.transloco.translate("auth.passwords_do_not_match"),
      );
      return;
    }
    this.authService
      .register({
        invitationCode: this.invitationCodeCtrl.value,
        username: this.usernameCtrl.value!,
        password: this.passwordCtrl.value!,
        email: this.emailCtrl.value,
      })
      .pipe(take(1))
      .subscribe(() => {
        void this.router.navigate(["/login"]);
      });
  }
}

class PasswordEntropyMatcher implements ErrorStateMatcher {
  private result?: generated.PasswordEntropyResult;

  constructor(result$: Observable<generated.PasswordEntropyResult>) {
    result$.subscribe((result) => {
      this.result = result;
    });
  }

  isErrorState(control: FormControl | null): boolean {
    return !!(control?.value && !(this.result?.valid ?? true));
  }
}

class ConfirmPasswordMatcher implements ErrorStateMatcher {
  constructor(private passwordControl: FormControl) {}

  isErrorState(control: FormControl | null): boolean {
    return !!(
      control?.value &&
      (control?.value || this.passwordControl.value) &&
      control?.value !== this.passwordControl.value
    );
  }
}
