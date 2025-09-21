import { Component, Injectable, inject } from "@angular/core";
import {
  MatSnackBar,
  MatSnackBarRef,
  MAT_SNACK_BAR_DATA,
} from "@angular/material/snack-bar";
import { AppModule } from "../app.module";

@Injectable({ providedIn: "root" })
export class ErrorsService {
  public readonly expiry = 1000 * 10;

  constructor(private snackBar: MatSnackBar) {}

  addError(error: string, expiry = this.expiry) {
    this.snackBar.openFromComponent(ErrorComponent, {
      duration: expiry,
      panelClass: ["snack-bar-error"],
      data: { error },
    });
  }
}

@Component({
  selector: "app-error",
  template: `
    <ng-container *transloco="let t">
      <label matSnackBarLabel>{{ error }}</label>
      <span matSnackBarActions>
        <button
          mat-mini-fab
          color="neutral"
          matSnackBarAction
          [matTooltip]="t('general.dismiss')"
          (click)="snackBarRef.dismissWithAction()"
        >
          <mat-icon>close</mat-icon>
        </button>
      </span>
    </ng-container>
  `,
  styles: `
    :host {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    ::ng-deep .mat-mdc-snack-bar-container .mat-mdc-snackbar-surface {
      background: var(--mat-form-field-error-text-color);
    }

    label {
      font-weight: bold;
      font-size: 16px;
    }

    button {
      margin-left: 20px;
      width: 30px;
      height: 30px;
    }
  `,
  imports: [AppModule],
  standalone: true,
})
export class ErrorComponent {
  snackBarRef = inject(MatSnackBarRef);

  error: string;

  constructor() {
    this.error = (inject(MAT_SNACK_BAR_DATA) as { error: string }).error;
  }
}
