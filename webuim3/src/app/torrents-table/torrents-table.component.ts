import { Component, Input } from '@angular/core';
import * as generated from '../graphql/generated';
import {DataSource} from "@angular/cdk/collections";
import {
  MatCell,
  MatCellDef,
  MatColumnDef,
  MatHeaderCell,
  MatHeaderCellDef,
  MatHeaderRow, MatHeaderRowDef, MatRow, MatRowDef,
  MatTable
} from "@angular/material/table";
import {MatCheckbox} from "@angular/material/checkbox";
import {MatIcon} from "@angular/material/icon";
import {MatChip, MatChipSet} from "@angular/material/chips";
import {HumanTimePipe} from "../pipes/human-time.pipe";
import {MatTooltip} from "@angular/material/tooltip";
import {TorrentContentComponent} from "../torrent-content/torrent-content.component";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {animate, state, style, transition, trigger} from "@angular/animations";
import {FilesizePipe} from "../pipes/filesize.pipe";

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
  ],
  templateUrl: './torrents-table.component.html',
  styleUrl: './torrents-table.component.scss',
  animations: [
    trigger('detailExpand', [
      state('collapsed,void', style({height: '0px', minHeight: '0'})),
      state('expanded', style({height: '*'})),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class TorrentsTableComponent {
  @Input() dataSource: DataSource<generated.TorrentContent>

  displayedColumns = [
    "select",
    "summary",
    "size",
    "publishedAt",
    "peers",
    "magnet",
  ];

  expandedTorrentContentId: string | undefined;
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
