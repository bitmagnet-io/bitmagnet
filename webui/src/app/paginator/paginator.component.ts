import {
  Component,
  EventEmitter,
  Input,
  numberAttribute,
  Output,
} from "@angular/core";
import { AppModule } from "../app.module";
import { IntEstimatePipe } from "../pipes/int-estimate.pipe";
import type { PageEvent } from "./paginator.types";

@Component({
  selector: "app-paginator",
  template: `
    <ng-container *transloco="let t">
      <div class="paginator">
        <mat-form-field class="field-items-per-page" subscriptSizing="dynamic">
          <mat-label>Items per page</mat-label>
          <mat-select
            [value]="pageSize"
            (valueChange)="pageSize = $event; page = 1; emitChange()"
          >
            @for (size of pageSizes; track size) {
              <mat-option [value]="size">
                {{ size }}
              </mat-option>
            }
          </mat-select>
        </mat-form-field>
        <p class="paginator-description">
          @if (hasTotalLength) {
            {{
              t("paginator.x_to_y_of_z", {
                x: (firstItemIndex | number),
                y: (lastItemIndex | number),
                z: totalLength ?? 0 | intEstimate: totalIsEstimate,
              })
            }}
          } @else {
            {{
              t("paginator.x_to_y", {
                x: (firstItemIndex | number),
                y: (lastItemIndex | number),
              })
            }}
          }
        </p>
        <div class="paginator-navigation">
          <button
            mat-icon-button
            [disabled]="!hasPreviousPage"
            (click)="page = 1; emitChange()"
            [matTooltip]="t('paginator.first_page')"
          >
            <mat-icon>first_page</mat-icon>
          </button>
          <button
            mat-icon-button
            [disabled]="!hasPreviousPage"
            (click)="page = page - 1; emitChange()"
            [matTooltip]="t('paginator.previous_page')"
          >
            <mat-icon>navigate_before</mat-icon>
          </button>
          <button
            mat-icon-button
            [disabled]="!actuallyHasNextPage"
            (click)="page = page + 1; emitChange()"
            [matTooltip]="t('paginator.next_page')"
          >
            <mat-icon>navigate_next</mat-icon>
          </button>
          @if (showLastPage) {
            <button
              mat-icon-button
              [disabled]="[null, page].includes(pageCount)"
              (click)="page = pageCount ?? 1; emitChange()"
              [matTooltip]="t('paginator.last_page')"
            >
              <mat-icon>last_page</mat-icon>
            </button>
          }
        </div>
      </div>
    </ng-container>
  `,
  styles: [
    `
      .paginator {
        > * {
          display: inline-block;
          vertical-align: middle;
        }
        p {
          margin: 0 20px;
        }
        .field-items-per-page {
          width: 140px;
        }
      }
    `,
  ],
  standalone: true,
  imports: [AppModule, IntEstimatePipe],
})
export class PaginatorComponent {
  @Input({ transform: numberAttribute }) page = 1;
  @Input({ transform: numberAttribute }) pageSize = 10;
  @Input() pageSizes: number[] = [10, 20, 50, 100];
  @Input({ transform: numberAttribute }) pageLength = 0;
  @Input() totalLength: number | null = null;
  @Input() totalIsEstimate = false;
  @Input() hasNextPage: boolean | null | undefined = null;
  @Input() showLastPage = false;

  @Output() paging = new EventEmitter<PageEvent>();

  get firstItemIndex() {
    return (this.page - 1) * this.pageSize + 1;
  }

  get lastItemIndex() {
    return (this.page - 1) * this.pageSize + this.pageLength;
  }

  get hasTotalLength() {
    return typeof this.totalLength === "number";
  }

  get hasPreviousPage() {
    return this.page > 1;
  }

  get pageCount(): number | null {
    if (typeof this.totalLength !== "number") {
      return null;
    }
    return Math.ceil(this.totalLength / this.pageSize);
  }

  get actuallyHasNextPage() {
    if (typeof this.hasNextPage === "boolean") {
      return this.hasNextPage;
    }
    if (typeof this.totalLength !== "number") {
      return false;
    }
    return this.page * this.pageSize < this.totalLength;
  }

  emitChange() {
    this.paging.emit({
      page: this.page,
      pageSize: this.pageSize,
    });
  }
}
