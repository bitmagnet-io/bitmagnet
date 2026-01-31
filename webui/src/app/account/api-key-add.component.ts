import { Component, inject } from "@angular/core";
import { MatDialogRef } from "@angular/material/dialog";
import { FormControl } from "@angular/forms";
import { Apollo } from "apollo-angular";
import { catchError, EMPTY, map, take } from "rxjs";
import { ErrorsService } from "../errors/errors.service";
import * as generated from "../graphql/generated";
import { AppModule } from "../app.module";
import {
  ChangeEvent,
  PermissionsEditComponent,
} from "../auth/permissions-edit.component";
import { RolesService } from "../auth/roles.service";
import { objectActionKey } from "../auth/util";

@Component({
  template: `
    <ng-container *transloco="let t">
      <mat-card>
        <mat-card-header>
          <h2>
            {{
              result ? t("account.api_key_created") : t("account.new_api_key")
            }}
          </h2>
        </mat-card-header>
        <mat-card-content>
          @if (!result) {
            <mat-form-field>
              <mat-label>{{ t("general.name") }}</mat-label>
              <input required matInput [formControl]="nameCtrl" />
            </mat-form-field>
            <app-permissions-edit
              [permissions]="(permissions$ | async)!"
              (change)="this.permissionChanges = $event"
            />
          } @else {
            <p>{{ t("account.api_key_take_note") }}</p>
            <p
              class="copy"
              [matTooltip]="t('general.copy_to_clipboard')"
              [cdkCopyToClipboard]="result.apiKey"
            >
              {{ result.apiKey }}
            </p>
          }
        </mat-card-content>
        <mat-card-footer>
          <mat-card-actions>
            @if (!result) {
              <button
                [disabled]="!nameCtrl.value"
                mat-stroked-button
                (click)="create()"
              >
                {{ t("general.confirm") }}
              </button>
            }
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
      h5 {
        margin-bottom: 10px;
      }
      p.copy {
        cursor: crosshair;
        text-decoration: underline;
        text-decoration-style: dotted;
        padding-bottom: 20px;
        margin-top: 0;
      }
    `,
  ],
  standalone: true,
  imports: [AppModule, PermissionsEditComponent],
})
export class APIKeyAddComponent {
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

  private dialogRef = inject(MatDialogRef<APIKeyAddComponent>);

  nameCtrl = new FormControl<string>("");

  permissionChanges?: ChangeEvent;

  result?: generated.CreateApiKeyResult;

  close() {
    this.dialogRef.close();
  }

  create() {
    this.apollo
      .mutate<
        generated.CreateApiKeyMutation,
        generated.CreateApiKeyMutationVariables
      >({
        mutation: generated.CreateApiKeyDocument,
        variables: {
          input: {
            name: this.nameCtrl.value!,
            permissions:
              this.permissionChanges?.permissions.flatMap((permission) =>
                permission.active && !permission.core
                  ? [permission.objectAction]
                  : [],
              ) ?? [],
          },
        },
      })
      .pipe(
        take(1),
        catchError((err, {}) => {
          this.errors.addError(`${err}`);
          return EMPTY;
        }),
        map((result) => {
          this.result = result.data?.self.createAPIKey;
        }),
      )
      .subscribe();
  }
}
