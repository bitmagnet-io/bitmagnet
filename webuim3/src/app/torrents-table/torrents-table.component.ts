import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {
  MatCell,
  MatCellDef,
  MatColumnDef,
  MatHeaderCell,
  MatHeaderCellDef,
  MatHeaderRow,
  MatHeaderRowDef,
  MatRow,
  MatRowDef,
  MatTable,
} from '@angular/material/table';
import { MatCheckbox } from '@angular/material/checkbox';
import { MatIcon } from '@angular/material/icon';
import { MatChip, MatChipAvatar, MatChipSet } from '@angular/material/chips';
import { MatTooltip } from '@angular/material/tooltip';
import {
  animate,
  state,
  style,
  transition,
  trigger,
} from '@angular/animations';
import { TranslocoDirective } from '@jsverse/transloco';
import { AsyncPipe } from '@angular/common';
import { MatProgressBar } from '@angular/material/progress-bar';
import { SelectionModel } from '@angular/cdk/collections';
import { FilesizePipe } from '../pipes/filesize.pipe';
import { TorrentContentComponent } from '../torrent-content/torrent-content.component';
import { HumanTimePipe } from '../pipes/human-time.pipe';
import * as generated from '../graphql/generated';
import { TorrentsSearchDatasource } from '../torrents-search/torrents-search.datasource';
import { contentTypeInfo } from '../taxonomy/content-types';
import {BehaviorSubject} from "rxjs";

@Component({
  selector: 'app-torrents-table',
  standalone: true,
  imports: [
    MatTable,
    MatColumnDef,
    MatHeaderCell,
    MatCheckbox,
    MatHeaderCellDef,
    MatCell,
    MatCellDef,
    MatIcon,
    MatChipSet,
    MatChip,
    HumanTimePipe,
    MatTooltip,
    MatHeaderRow,
    MatRow,
    MatRowDef,
    MatHeaderRowDef,
    TorrentContentComponent,
    FilesizePipe,
    TranslocoDirective,
    AsyncPipe,
    MatProgressBar,
    MatChipAvatar,
  ],
  templateUrl: './torrents-table.component.html',
  styleUrl: './torrents-table.component.scss',
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
export class TorrentsTableComponent implements OnInit {
  contentTypeInfo = contentTypeInfo;

  @Input() dataSource: TorrentsSearchDatasource;
  @Input() selection: SelectionModel<string>;
  @Input() displayedColumns: readonly Column[] = allColumns;

  @Output() updated = new EventEmitter<string>();

  expandedTorrentContentId: string | undefined;

  items = Array<generated.TorrentContent>();

  ngOnInit() {
    this.dataSource.items$.subscribe((items) => {
      this.items = items;
    });
  }

  /** Whether the number of selected elements matches the total number of rows. */
  isAllSelected() {
    return this.items.every((i) => this.selection.isSelected(i.infoHash));
  }

  /** Selects all rows if they are not all selected; otherwise clear selection. */
  toggleAllRows() {
    if (this.isAllSelected()) {
      this.selection.clear();
      return;
    }
    this.selection.select(...this.items.map((i) => i.infoHash));
  }

  toggleTorrentContentId(id: string) {
    if (this.expandedTorrentContentId === id) {
      this.expandedTorrentContentId = undefined;
    } else {
      this.expandedTorrentContentId = id;
    }
  }

  /**
   * Workaround for untyped table cell definitions
   */
  item(item: generated.TorrentContent): generated.TorrentContent {
    return item;
  }
}

export const allColumns = [
  'select',
  'summary',
  'size',
  'publishedAt',
  'peers',
  'magnet',
] as const;

export const compactColumns = ['select', 'summary', 'size', 'magnet'] as const;

export type Column = (typeof allColumns)[number];
