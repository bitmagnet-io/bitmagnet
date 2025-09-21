import { Component, inject } from "@angular/core";
import { AppModule } from "../../app.module";
import { MatDialogRef } from "@angular/material/dialog";
import { FormControl } from "@angular/forms";
import { Apollo } from "apollo-angular";
import { ErrorsService } from "../../errors/errors.service";
import * as generated from "../../graphql/generated";
import { catchError, EMPTY, map, Observable, take } from "rxjs";
import {
  ChangeEvent,
  PermissionsEditComponent,
} from "../../auth/permissions-edit.component";
import { RolesService } from "../../auth/roles.service";
import { objectActionKey } from "../../auth/util";

@Component({
  template: `
    <ng-container *transloco="let t">
      <mat-card>
        <mat-card-header>
          <h2>{{ t("auth.new_role") }}</h2>
        </mat-card-header>
        <mat-card-content>
          <mat-form-field>
            <mat-label>{{ t("general.name") }}</mat-label>
            <input required matInput [formControl]="nameCtrl" />
          </mat-form-field>
          <app-permissions-edit
            [permissions]="(permissions$ | async)!"
            (change)="this.permissionChanges = $event"
          />
        </mat-card-content>
        <mat-card-footer>
          <mat-card-actions>
            <button
              [disabled]="!nameCtrl.value"
              mat-stroked-button
              (click)="create()"
            >
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
      app-permissions-edit {
        max-height: 300px;
        overflow-y: scroll;
        overflow-x: hidden;
      }
    `,
  ],
  standalone: true,
  imports: [AppModule, PermissionsEditComponent],
})
export class RoleAddComponent {
  private apollo = inject(Apollo);
  private errors = inject(ErrorsService);
  private rolesService = inject(RolesService);

  permissions$ = this.rolesService.objectActions$.pipe(
    take(1),
    map((objectActions) =>
      objectActions.map((objectAction) => ({
        key: objectActionKey(objectAction),
        objectAction,
        active: false,
        core: false,
      })),
    ),
  );

  private dialogRef = inject(MatDialogRef<RoleAddComponent>);

  nameCtrl = new FormControl<string>("");

  permissionChanges?: ChangeEvent;

  close() {
    this.dialogRef.close();
  }

  create() {
    this.apollo
      .mutate<generated.PutRoleMutation, generated.PutRoleMutationVariables>({
        mutation: generated.PutRoleDocument,
        variables: {
          role: this.nameCtrl.value!,
          objectActions:
            this.permissionChanges?.permissions.flatMap((permission) =>
              permission.active && !permission.core
                ? [permission.objectAction]
                : [],
            ) ?? [],
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
