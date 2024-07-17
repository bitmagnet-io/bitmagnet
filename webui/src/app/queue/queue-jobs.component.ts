import {Component, inject, OnInit, ViewChild} from '@angular/core';
import {Apollo} from "apollo-angular";
import {GraphQLModule} from "../graphql/graphql.module";
import {ErrorsService} from "../errors/errors.service";
import {AsyncPipe} from "@angular/common";
import {FilesizePipe} from "../pipes/filesize.pipe";
import {
  MatCell,
  MatCellDef,
  MatColumnDef,
  MatHeaderCell, MatHeaderCellDef,
  MatHeaderRow,
  MatHeaderRowDef,
  MatRow, MatRowDef, MatTable
} from "@angular/material/table";
import {MatCheckbox} from "@angular/material/checkbox";
import {MatIcon} from "@angular/material/icon";
import {MatProgressBar} from "@angular/material/progress-bar";
import {MatTooltip} from "@angular/material/tooltip";
import {TimeAgoPipe} from "../dates/time-ago.pipe";
import {TorrentChipsComponent} from "../torrent-chips/torrent-chips.component";
import {TorrentContentComponent} from "../torrent-content/torrent-content.component";
import {TranslocoDirective, TranslocoService} from "@jsverse/transloco";
import {QueueJobsController} from "./queue-jobs.controller";
import {QueueJobsDatasource} from "./queue-jobs.datasource";
import {BehaviorSubject} from "rxjs";
import * as generated from "../graphql/generated";
import {SelectionModel} from "@angular/cdk/collections";
import {Column} from "../torrents-table/torrents-table.component";
import {animate, state, style, transition, trigger} from "@angular/animations";

@Component({
  selector: 'app-queue-jobs',
  standalone: true,
  imports: [
    GraphQLModule,
    AsyncPipe,
    FilesizePipe,
    MatCell,
    MatCellDef,
    MatCheckbox,
    MatColumnDef,
    MatHeaderCell,
    MatHeaderRow,
    MatHeaderRowDef,
    MatIcon,
    MatProgressBar,
    MatRow,
    MatRowDef,
    MatTable,
    MatTooltip,
    TimeAgoPipe,
    TorrentChipsComponent,
    TorrentContentComponent,
    TranslocoDirective,
    MatHeaderCellDef,
  ],
  templateUrl: './queue-jobs.component.html',
  styleUrl: './queue-jobs.component.scss',
  animations: [
    trigger('detailExpand', [
      state('collapsed,void', style({ height: '0px', minHeight: '0' })),
      state('expanded', style({ height: '*' })),
      transition(
        'expanded <=> collapsed',
        animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)'),
      ),
    ]),
  ],
})
export class QueueJobsComponent implements OnInit {
  private apollo = inject(Apollo);
  private errorsService = inject(ErrorsService);
  protected transloco = inject(TranslocoService);
  protected controller = new QueueJobsController()
  protected dataSource = new QueueJobsDatasource(this.apollo, this.errorsService, this.controller.variables$)

  expandedId = new BehaviorSubject<string | null>(null);

  items = Array<generated.QueueJob>();

  @ViewChild('selection') selection: SelectionModel<string>;

  displayedColumns = ["id", "queue", "status", "error", "createdAt"];

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
}
