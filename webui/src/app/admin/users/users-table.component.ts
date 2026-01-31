import { Component } from "@angular/core";
import { BehaviorSubject } from "rxjs";
import {
  animate,
  state,
  style,
  transition,
  trigger,
} from "@angular/animations";
import { UsersDatasource } from "../../auth/users.datasource";
import { AppModule } from "../../app.module";
import * as generated from "../../graphql/generated";
import { TimeAgoPipe } from "../../pipes/time-ago.pipe";
import { PaginatorComponent } from "../../paginator/paginator.component";

@Component({
  selector: "app-users-table",
  template: `
    <ng-container *transloco="let t">
      <table
        mat-table
        [dataSource]="dataSource"
        [multiTemplateDataRows]="true"
        class="table-results"
      >
        <ng-container matColumnDef="id">
          <th mat-header-cell *matHeaderCellDef>ID</th>
          <td
            mat-cell
            *matCellDef="let i"
            (click)="toggleUser(i.id); $event.stopPropagation()"
          >
            {{ item(i).id }}
          </td>
        </ng-container>

        <ng-container matColumnDef="username">
          <th mat-header-cell *matHeaderCellDef>
            {{ t("auth.username") }}
          </th>
          <td
            mat-cell
            *matCellDef="let i"
            (click)="toggleUser(i.id); $event.stopPropagation()"
          >
            {{ item(i).username }}
          </td>
        </ng-container>

        <ng-container matColumnDef="role">
          <th mat-header-cell *matHeaderCellDef>
            {{ t("auth.role") }}
          </th>
          <td
            mat-cell
            *matCellDef="let i"
            (click)="toggleUser(i.id); $event.stopPropagation()"
          >
            {{ item(i).role }}
          </td>
        </ng-container>

        <ng-container matColumnDef="email">
          <th mat-header-cell *matHeaderCellDef>
            {{ t("general.email") }}
          </th>
          <td
            mat-cell
            *matCellDef="let i"
            (click)="toggleUser(i.id); $event.stopPropagation()"
          >
            {{ item(i).email }}
          </td>
        </ng-container>

        <ng-container matColumnDef="createdAt">
          <th mat-header-cell *matHeaderCellDef>
            {{ t("general.created_at") }}
          </th>
          <td
            mat-cell
            *matCellDef="let i"
            [matTooltip]="item(i).createdAt"
            (click)="toggleUser(i.id); $event.stopPropagation()"
          >
            {{ item(i).createdAt | timeAgo }}
          </td>
        </ng-container>

        <ng-container matColumnDef="lastLoginAt">
          <th mat-header-cell *matHeaderCellDef>
            {{ t("auth.last_login_at") }}
          </th>
          <td
            mat-cell
            *matCellDef="let i"
            [matTooltip]="i.lastLoginAt ?? t('general.never')"
            (click)="toggleUser(i.id); $event.stopPropagation()"
          >
            @if (i.lastLoginAt; as lastLoginAt) {
              {{ lastLoginAt | timeAgo }}
            } @else {
              {{ t("general.never") }}
            }
          </td>
        </ng-container>

        <ng-container matColumnDef="expandedDetail">
          <td
            mat-cell
            *matCellDef="let i"
            [attr.colspan]="displayedColumns.length"
          >
            <div
              class="item-detail"
              [@detailExpand]="
                expandedId.getValue() === i.id ? 'expanded' : 'collapsed'
              "
            >
              Hello
            </div>
          </td>
        </ng-container>

        <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
        <tr
          mat-row
          *matRowDef="let i; columns: displayedColumns"
          [class]="
            'summary-row ' +
            (i.id === expandedId.getValue() ? 'expanded' : 'collapsed')
          "
        ></tr>
        <tr
          mat-row
          *matRowDef="let i; columns: ['expandedDetail']"
          [class]="
            'expanded-detail-row ' +
            (i.id === expandedId.getValue() ? 'expanded' : 'collapsed')
          "
        ></tr>
      </table>
      <app-paginator
        (paging)="dataSource.handlePagination($event)"
        [page]="dataSource.page"
        [pageSize]="dataSource.limit"
        [pageLength]="(dataSource.users$ | async)?.length ?? 0"
        [totalLength]="(dataSource.result$ | async)?.totalCount ?? 0"
        [totalIsEstimate]="false"
        [showLastPage]="true"
      />
    </ng-container>
  `,
  styles: [
    `
      tr:not(.expanded-detail-row) td {
        cursor: pointer;
      }

      tr.expanded-detail-row {
        height: 0;
      }

      app-paginator {
        margin-top: 10px;
        float: right;
      }
    `,
  ],
  animations: [
    trigger("detailExpand", [
      state("collapsed,void", style({ height: "0px", minHeight: "0" })),
      state("expanded", style({ height: "*" })),
      transition(
        "expanded <=> collapsed",
        animate("225ms cubic-bezier(0.4, 0.0, 0.2, 1)"),
      ),
    ]),
  ],
  standalone: true,
  imports: [AppModule, PaginatorComponent, TimeAgoPipe],
})
export class UsersTableComponent {
  dataSource = new UsersDatasource();

  displayedColumns = [
    "id",
    "username",
    "role",
    "email",
    "createdAt",
    "lastLoginAt",
  ];

  expandedId = new BehaviorSubject<number | null>(null);

  item(item: unknown): generated.User {
    return item as generated.User;
  }

  toggleUser(id: number) {
    const current = this.expandedId.getValue();
    this.expandedId.next(current === id ? null : id);
  }
}
