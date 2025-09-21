import { Component, inject, Input, OnInit } from "@angular/core";
import { InvitationsDatasource } from "../../auth/invitations.datasource";
import { AppModule } from "../../app.module";
import * as generated from "../../graphql/generated";
import { TimeAgoPipe } from "../../pipes/time-ago.pipe";
import { Apollo } from "apollo-angular";
import { Observable } from "rxjs";
import { PaginatorComponent } from "../../paginator/paginator.component";

@Component({
  selector: "app-invitations-table",
  template: `
    <ng-container *transloco="let t">
      <table
        mat-table
        [dataSource]="dataSource"
        [multiTemplateDataRows]="true"
        class="table-results"
      >
        <ng-container matColumnDef="code">
          <th mat-header-cell *matHeaderCellDef>{{ t("general.code") }}</th>
          <td mat-cell *matCellDef="let i">
            <span
              class="copy"
              [matTooltip]="t('general.copy_to_clipboard')"
              [cdkCopyToClipboard]="item(i).code"
              >{{ item(i).code }}</span
            >
          </td>
        </ng-container>

        <ng-container matColumnDef="role">
          <th mat-header-cell *matHeaderCellDef>
            {{ t("auth.role") }}
          </th>
          <td mat-cell *matCellDef="let i">
            {{ item(i).role }}
          </td>
        </ng-container>

        <ng-container matColumnDef="email">
          <th mat-header-cell *matHeaderCellDef>
            {{ t("general.email") }}
          </th>
          <td mat-cell *matCellDef="let i">
            {{ item(i).email }}
          </td>
        </ng-container>

        <ng-container matColumnDef="createdBy">
          <th mat-header-cell *matHeaderCellDef>{{ t("auth.created_by") }}</th>
          <td mat-cell *matCellDef="let i">
            {{ item(i).createdBy?.username || t("general.system") }}
          </td>
        </ng-container>

        <ng-container matColumnDef="claimedBy">
          <th mat-header-cell *matHeaderCellDef>{{ t("auth.claimed_by") }}</th>
          <td mat-cell *matCellDef="let i">
            {{ item(i).claimedBy?.username || t("auth.unclaimed") }}
          </td>
        </ng-container>

        <ng-container matColumnDef="createdAt">
          <th mat-header-cell *matHeaderCellDef>
            {{ t("general.created_at") }}
          </th>
          <td mat-cell *matCellDef="let i">
            <span [matTooltip]="item(i).createdAt">{{
              item(i).createdAt | timeAgo
            }}</span>
          </td>
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

        <ng-container matColumnDef="actions">
          <th mat-header-cell *matHeaderCellDef style="text-align: center">
            {{ t("general.actions") }}
          </th>
          <td mat-cell *matCellDef="let i" style="text-align: center">
            <button mat-stroked-button (click)="delete(i.code)">
              {{ t("general.delete") }}
            </button>
          </td>
        </ng-container>

        <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
        <tr mat-row *matRowDef="let i; columns: displayedColumns"></tr>
      </table>
      <app-paginator
        (paging)="dataSource.handlePagination($event)"
        [page]="dataSource.page"
        [pageSize]="dataSource.limit"
        [pageLength]="(dataSource.invitations$ | async)?.length ?? 0"
        [totalLength]="(dataSource.result$ | async)?.totalCount ?? 0"
        [totalIsEstimate]="false"
        [showLastPage]="true"
      />
    </ng-container>
  `,
  styles: [
    `
      span.copy {
        cursor: crosshair;
        text-decoration: underline;
        text-decoration-style: dotted;
      }
      app-paginator {
        margin-top: 10px;
        float: right;
      }
    `,
  ],
  standalone: true,
  imports: [AppModule, PaginatorComponent, TimeAgoPipe],
})
export class InvitationsTableComponent implements OnInit {
  private apollo = inject(Apollo);
  dataSource = new InvitationsDatasource();

  @Input() update: Observable<void>;

  displayedColumns = [
    "code",
    "role",
    "email",
    "createdBy",
    "claimedBy",
    "createdAt",
    "expiresAt",
    "actions",
  ];

  ngOnInit(): void {
    this.update.subscribe(() => this.dataSource.refresh());
  }

  item(item: any): generated.Invitation {
    return item;
  }

  delete(code: string) {
    this.apollo
      .mutate<
        generated.InvitationDeleteMutation,
        generated.InvitationDeleteMutationVariables
      >({
        mutation: generated.InvitationDeleteDocument,
        variables: {
          code,
        },
      })
      .subscribe(() => this.dataSource.refresh());
  }
}
