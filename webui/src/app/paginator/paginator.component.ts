import {
  Component,
  EventEmitter,
  Input,
  numberAttribute,
  Output,
} from "@angular/core";
import { AppModule } from "../app.module";
import type { PageEvent } from "./paginator.types";

@Component({
  selector: "app-paginator",
  templateUrl: "./paginator.component.html",
  standalone: true,
  styleUrls: ["./paginator.component.scss"],
  imports: [AppModule],
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
