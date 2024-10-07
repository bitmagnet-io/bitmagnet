import { Component, inject, Input, OnInit } from '@angular/core';
import { Apollo } from 'apollo-angular';
import { AsyncPipe, DecimalPipe, SlicePipe } from '@angular/common';
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
import { MatProgressBar } from '@angular/material/progress-bar';
import { TranslocoDirective, TranslocoService } from '@jsverse/transloco';
import { TimeAgoPipe } from '../dates/time-ago.pipe';
import * as generated from '../graphql/generated';
import { ErrorsService } from '../errors/errors.service';
import { PaginatorComponent } from '../paginator/paginator.component';
import { QueueJobsTableComponent } from '../queue/queue-jobs-table.component';
import { FilesizePipe } from '../pipes/filesize.pipe';
import {
  ITorrentFilesDatasource,
  TorrentFilesDatasource,
  TorrentFilesSingleDatasource,
} from './torrent-files.datasource';
import {
  TorrentFilesController,
  TorrentFilesControls,
} from './torrent-files.controller';

@Component({
  selector: 'app-torrent-files-table',
  standalone: true,
  imports: [
    AsyncPipe,
    DecimalPipe,
    MatCell,
    MatCellDef,
    MatColumnDef,
    MatHeaderCell,
    MatHeaderRow,
    MatHeaderRowDef,
    MatProgressBar,
    MatRow,
    MatRowDef,
    MatTable,
    SlicePipe,
    TimeAgoPipe,
    TranslocoDirective,
    MatHeaderCellDef,
    PaginatorComponent,
    QueueJobsTableComponent,
    FilesizePipe,
  ],
  templateUrl: './torrent-files-table.component.html',
  styleUrl: './torrent-files-table.component.scss',
})
export class TorrentFilesTableComponent implements OnInit {
  private apollo = inject(Apollo);
  private errorsService = inject(ErrorsService);
  protected transloco = inject(TranslocoService);

  @Input() torrent: generated.Torrent;

  protected controller: TorrentFilesController;
  protected dataSource: ITorrentFilesDatasource;

  protected displayedColumns = ['index', 'path', 'type', 'size'];

  protected controls: TorrentFilesControls;

  ngOnInit() {
    this.controller = new TorrentFilesController(this.torrent.infoHash);
    this.dataSource =
      this.torrent.filesStatus === 'single'
        ? new TorrentFilesSingleDatasource(this.torrent)
        : new TorrentFilesDatasource(
            this.apollo,
            this.errorsService,
            this.controller.variables$,
          );
    this.controller.controls$.subscribe((ctrl) => {
      this.controls = ctrl;
    });
  }

  /**
   * Workaround for untyped table cell definitions
   */
  item(item: generated.TorrentFile): generated.TorrentFile {
    return item;
  }
}
