import { Injectable } from "@angular/core";
import { MatPaginatorIntl } from "@angular/material/paginator";
import { formatNumber as _formatNumber } from "@angular/common";

@Injectable()
export class MatPaginatorIntlCustom extends MatPaginatorIntl {
  override getRangeLabel = (page: number, pageSize: number, length: number) => {
    return [
      formatNumber(page * pageSize + 1),
      "-",
      Math.min((page + 1) * pageSize, length),
      "of ≤",
      formatNumber(length),
    ].join(" ");
  };
}

const formatNumber = (value: number) => _formatNumber(value, "en-US");
