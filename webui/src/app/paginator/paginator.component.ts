import {
  Component,
  EventEmitter,
  Input,
  numberAttribute,
  Output,
} from "@angular/core";
import type { PageEvent } from "./paginator.types";

@Component({
  selector: "app-paginator",
  templateUrl: "./paginator.component.html",
  styleUrls: ["./paginator.component.scss"],
})
export class PaginatorComponent {
  @Input({ transform: numberAttribute }) pageIndex = 0;
  @Input({ transform: numberAttribute }) pageSize = 10;
  @Input() pageSizes: number[] = [10, 20, 50, 100];
  @Input({ transform: numberAttribute }) pageLength = 0;
  @Input() totalLength: number | null = null;
  @Input() totalLessThanOrEqual = false;

  @Output() page = new EventEmitter<PageEvent>();

  get firstItemIndex() {
    return this.pageIndex * this.pageSize + 1;
  }

  get lastItemIndex() {
    return this.pageIndex * this.pageSize + this.pageLength;
  }

  get hasTotalLength() {
    return typeof this.totalLength === "number";
  }

  get hasPreviousPage() {
    return this.pageIndex > 0;
  }

  get hasNextPage() {
    return this.firstItemIndex + this.pageSize <= this.totalLength!;
  }

  emitChange() {
    this.page.emit({
      pageIndex: this.pageIndex,
      pageSize: this.pageSize,
    });
  }
}
