import {
  Component,
  EventEmitter,
  inject,
  Input,
  OnInit,
  Output,
} from "@angular/core";
import { SelectionModel } from "@angular/cdk/collections";
import { TranslocoService } from "@jsverse/transloco";
import { BehaviorSubject } from "rxjs";
import {
  animate,
  state,
  style,
  transition,
  trigger,
} from "@angular/animations";
import { TimeAgoPipe } from "../../pipes/time-ago.pipe";
import * as generated from "../../graphql/generated";
import { AppModule } from "../../app.module";
import { QueueJobsDatasource } from "./queue-jobs.datasource";

@Component({
  selector: "app-queue-jobs-table",
  standalone: true,
  imports: [AppModule, TimeAgoPipe],
  templateUrl: "./queue-jobs-table.component.html",
  styleUrl: "./queue-jobs-table.component.scss",
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
})
export class QueueJobsTableComponent implements OnInit {
  protected transloco = inject(TranslocoService);

  @Input() dataSource: QueueJobsDatasource;
  @Input() selection: SelectionModel<string>;
  @Input() displayedColumns = allColumns;

  @Output() updated = new EventEmitter<string>();

  expandedId = new BehaviorSubject<string | null>(null);

  items = Array<generated.QueueJob>();

  ngOnInit() {
    this.dataSource.items$.subscribe((items) => {
      this.items = items;
      if (items.length) {
        const expandedId = this.expandedId.getValue();
        if (expandedId && !items.some(({ id }) => id === expandedId)) {
          this.expandedId.next(null);
        }
      }
    });
  }

  /** Whether the number of selected elements matches the total number of rows. */
  isAllSelected() {
    return this.items.every((i) => this.selection.isSelected(i.id));
  }

  /** Selects all rows if they are not all selected; otherwise clear selection. */
  toggleAllRows() {
    if (this.isAllSelected()) {
      this.selection.clear();
      return;
    }
    this.selection.select(...this.items.map((i) => i.id));
  }

  toggleQueueJobId(id: string) {
    if (this.expandedId.getValue() === id) {
      this.expandedId.next(null);
    } else {
      this.expandedId.next(id);
    }
  }

  /**
   * Workaround for untyped table cell definitions
   */
  item(item: generated.QueueJob): generated.QueueJob {
    return item;
  }

  beautifyPayload(payload: string): string {
    try {
      return JSON.stringify(JSON.parse(payload), null, 2);
      //eslint-disable-next-line  @typescript-eslint/no-unused-vars
    } catch (e) {
      return payload;
    }
  }
}

const allColumns = [
  "id",
  "queue",
  "priority",
  "status",
  "error",
  "createdAt",
  "ranAt",
];
