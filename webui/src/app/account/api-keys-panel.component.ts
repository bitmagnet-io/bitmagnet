import { Component, inject } from "@angular/core";
import { AppModule } from "../app.module";
import { MatDialog } from "@angular/material/dialog";
import { Apollo } from "apollo-angular";
import * as generated from "../graphql/generated";
import { TimeAgoPipe } from "../pipes/time-ago.pipe";
import { APIKeyAddComponent } from "./api-key-add.component";
import { map, Observable } from "rxjs";

@Component({
  template: `
    <ng-container *transloco="let t">
      <app-document-title [parts]="[t('routes.api_keys')]" />
      <mat-card>
        <mat-card-header>
          <mat-toolbar>
            <h2><mat-icon>key</mat-icon>{{ t("routes.api_keys") }}</h2>
            <button mat-stroked-button (click)="newAPIKey()">
              <mat-icon>add</mat-icon>{{ t("account.new_api_key") }}
            </button>
          </mat-toolbar>
        </mat-card-header>
        <mat-card-content>
          <table mat-table [dataSource]="apiKeys$">
            <ng-container matColumnDef="name">
              <th mat-header-cell *matHeaderCellDef>{{ t("general.name") }}</th>
              <td mat-cell *matCellDef="let element">{{ element.name }}</td>
            </ng-container>
            <ng-container matColumnDef="expiresAt">
              <th mat-header-cell *matHeaderCellDef>
                {{ t("auth.expires_at") }}
              </th>
              <td mat-cell *matCellDef="let i">
                @if (i.expiresAt; as expiresAt) {
                  {{ expiresAt | timeAgo }}
                } @else {
                  {{ t("general.never") }}
                }
              </td>
            </ng-container>

            <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
            <tr mat-row *matRowDef="let row; columns: displayedColumns"></tr>
          </table>
        </mat-card-content>
      </mat-card>
    </ng-container>
  `,
  styles: [
    `
      mat-toolbar {
        h2 mat-icon {
          position: relative;
          top: 3px;
          margin-right: 14px;
          margin-left: 20px;
        }

        justify-content: space-between;
      }

      mat-expansion-panel {
        mat-divider {
          margin-bottom: 20px;
        }
      }
    `,
  ],

  standalone: true,
  imports: [AppModule, TimeAgoPipe],
})
export class APIKeysPanelComponent {
  private apollo = inject(Apollo);
  readonly dialog = inject(MatDialog);

  apiKeys$: Observable<generated.ApiKeyFragment[]>;

  displayedColumns = ["name", "expiresAt"];

  constructor() {
    this.refresh();
  }

  refresh() {
    this.apiKeys$ = this.apollo
      .query<generated.ApiKeysQuery, generated.ApiKeysQueryVariables>({
        query: generated.ApiKeysDocument,
        fetchPolicy: "no-cache",
      })
      .pipe(map((result) => result.data.self.apiKeys));
  }

  newAPIKey() {
    const ref = this.dialog.open(APIKeyAddComponent);
    ref.afterClosed().subscribe(() => {
      this.refresh();
    });
  }
}
