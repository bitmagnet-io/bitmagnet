import {
  Component,
  EventEmitter,
  Input,
  numberAttribute,
  Output,
} from "@angular/core";
import type { PageEvent } from "./paginator.types";
import {MatFormField, MatLabel} from "@angular/material/form-field";
import {MatOption, MatSelect} from "@angular/material/select";
import {DecimalPipe} from "@angular/common";
import {MatIconButton} from "@angular/material/button";
import {MatTooltip} from "@angular/material/tooltip";
import {MatIcon} from "@angular/material/icon";

@Component({
  selector: "app-paginator",
  templateUrl: "./paginator.component.html",
  standalone: true,
  styleUrls: ["./paginator.component.scss"],
  imports: [
    MatFormField,
    MatSelect,
    MatOption,
    DecimalPipe,
    MatIconButton,
    MatTooltip,
    MatIcon,
    MatLabel
  ]
})
export class PaginatorComponent {
  @Input({ transform: numberAttribute }) page = 1;
  @Input({ transform: numberAttribute }) pageSize = 10;
  @Input() pageSizes: number[] = [10, 20, 50, 100];
  @Input({ transform: numberAttribute }) pageLength = 0;
  @Input() totalLength: number | null = null;
  @Input() totalIsEstimate = false;
  @Input() hasNextPage: boolean | null | undefined = null;

  @Output() paging = new EventEmitter<PageEvent>();

  get firstItemIndex() {
    return (this.page-1) * this.pageSize + 1;
  }

  get lastItemIndex() {
    return (this.page-1) * this.pageSize + this.pageLength;
  }

  get hasTotalLength() {
    return typeof this.totalLength === "number";
  }

  get hasPreviousPage() {
    return this.page > 1;
  }

  emitChange() {
    this.paging.emit({
      page: this.page,
      pageSize: this.pageSize,
    });
  }
}
