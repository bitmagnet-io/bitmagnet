import {Component, inject, Input, OnInit} from '@angular/core';
import {GraphQLService} from "../graphql/graphql.service";
import {ErrorsService} from "../errors/errors.service";
import * as generated from "../graphql/generated";
import {TorrentFilesController, TorrentFilesControls} from "./torrent-files.controller";
import {
  ITorrentFilesDatasource,
  TorrentFilesDatasource,
  TorrentFilesSingleDatasource
} from "./torrent-files.datasource";
import {Apollo} from "apollo-angular";
import {AsyncPipe, DecimalPipe, SlicePipe} from "@angular/common";
import {
  MatCell,
  MatCellDef,
  MatColumnDef,
  MatHeaderCell, MatHeaderCellDef,
  MatHeaderRow,
  MatHeaderRowDef,
  MatRow, MatRowDef, MatTable
} from "@angular/material/table";
import {MatProgressBar} from "@angular/material/progress-bar";
import {TimeAgoPipe} from "../dates/time-ago.pipe";
import {TranslocoDirective, TranslocoService} from "@jsverse/transloco";
import {PaginatorComponent} from "../paginator/paginator.component";
import {QueueJobsTableComponent} from "../queue/queue-jobs-table.component";
import {QueueJobsControls} from "../queue/queue-jobs.controller";
import {FilesizePipe} from "../pipes/filesize.pipe";

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
    FilesizePipe
  ],
  templateUrl: './torrent-files-table.component.html',
  styleUrl: './torrent-files-table.component.scss'
})
export class TorrentFilesTableComponent implements  OnInit {
  private apollo = inject(Apollo)
  private errorsService = inject(ErrorsService);
  protected transloco = inject(TranslocoService);

  @Input() torrent: generated.Torrent;

  protected controller: TorrentFilesController
  protected dataSource: ITorrentFilesDatasource

  protected displayedColumns = ["index", "path", "type", "size"]

  protected controls: TorrentFilesControls;

  ngOnInit() {
    this.controller = new TorrentFilesController(this.torrent.infoHash)
    this.dataSource = this.torrent.filesStatus === "single" ? new TorrentFilesSingleDatasource(this.torrent) : new TorrentFilesDatasource(this.apollo, this.errorsService, this.controller.variables$)
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
