import { Component, inject } from "@angular/core";
import { MatDialogRef } from "@angular/material/dialog";
import { FormControl } from "@angular/forms";
import { Apollo } from "apollo-angular";
import { catchError, EMPTY, take } from "rxjs";
import { AppModule } from "../../app.module";
import { RolesService } from "../../auth/roles.service";
import { ErrorsService } from "../../errors/errors.service";
import * as generated from "../../graphql/generated";

@Component({
  template: `
    <ng-container *transloco="let t">
      <mat-card>
        <mat-card-header>
          <h2>{{ t("auth.new_invitation") }}</h2>
        </mat-card-header>
        <mat-card-content>
          <mat-form-field>
            <mat-label>{{ t("auth.role") }}</mat-label>
            <mat-select [formControl]="roleCtrl">
              @for (role of roles$ | async; track role.name) {
                <mat-option [value]="role.name">{{ role.name }}</mat-option>
              }
            </mat-select>
          </mat-form-field>
          <mat-form-field>
            <mat-label>{{ t("auth.expiry") }}</mat-label>
            <mat-select [formControl]="expiryCtrl">
              <mat-option value="P1D">{{
                t("admin.interval.days_1")
              }}</mat-option>
              <mat-option value="P1W">{{
                t("admin.interval.weeks_1")
              }}</mat-option>
            </mat-select>
          </mat-form-field>
        </mat-card-content>
        <mat-card-footer>
          <mat-card-actions>
            <button mat-stroked-button (click)="create()">
              {{ t("general.confirm") }}
            </button>
            <button mat-stroked-button (click)="close()">
              {{ t("general.dismiss") }}
            </button>
          </mat-card-actions>
        </mat-card-footer>
      </mat-card>
    </ng-container>
  `,
  styles: [
    `
      mat-card-content {
        display: flex;
        flex-direction: column;
        gap: 10px;
      }
      mat-card-actions {
        gap: 10px;
      }
    `,
  ],
  standalone: true,
  imports: [AppModule],
})
export class InvitationAddComponent {
  private apollo = inject(Apollo);
  private roles = inject(RolesService);
  private errors = inject(ErrorsService);
  private dialogRef = inject(MatDialogRef<InvitationAddComponent>);

  roles$ = this.roles.listRoles();

  roleCtrl = new FormControl<string>("user");
  expiryCtrl = new FormControl<string>("P1W");

  close() {
    this.dialogRef.close();
  }

  create() {
    this.apollo
      .mutate<generated.InviteMutation, generated.InviteMutationVariables>({
        mutation: generated.InviteDocument,
        variables: {
          input: {
            role: this.roleCtrl.value,
            expiry: this.expiryCtrl.value,
          },
        },
      })
      .pipe(
        take(1),
        catchError((err, {}) => {
          this.errors.addError(`${err}`);
          return EMPTY;
        }),
      )
      .subscribe(() => this.close());
  }
}
