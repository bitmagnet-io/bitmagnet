import { Component, EventEmitter, inject, Input, Output } from "@angular/core";
import { catchError, EMPTY, tap } from "rxjs";
import { NgOptimizedImage } from "@angular/common";
import { TranslocoService } from "@jsverse/transloco";
import { FilesizePipe } from "../pipes/filesize.pipe";
import * as generated from "../graphql/generated";
import { GraphQLService } from "../graphql/graphql.service";
import { ErrorsService } from "../errors/errors.service";
import { BreakpointsService } from "../layout/breakpoints.service";
import { TimeAgoPipe } from "../pipes/time-ago.pipe";
import { AppModule } from "../app.module";
import { TorrentFilesTableComponent } from "./torrent-files-table.component";
import { TorrentEditTagsComponent } from "./torrent-edit-tags.component";
import {
  TorrentTab,
  torrentTabNames,
  TorrentTabSelection,
} from "./torrents-search.controller";
import { TorrentReprocessComponent } from "./torrent-reprocess.component";

@Component({
  selector: "app-torrent-content",
  templateUrl: "./torrent-content.component.html",
  styleUrl: "./torrent-content.component.scss",
  standalone: true,
  imports: [
    AppModule,
    FilesizePipe,
    NgOptimizedImage,
    TimeAgoPipe,
    TorrentEditTagsComponent,
    TorrentFilesTableComponent,
    TorrentReprocessComponent,
  ],
})
export class TorrentContentComponent {
  breakpoints = inject(BreakpointsService);

  @Input() torrentContent: generated.TorrentContent;
  @Input() heading = true;
  @Input() size = true;
  @Input() peers = true;
  @Input() published = true;

  @Output() updated = new EventEmitter<null>();
  @Output() tabSelected = new EventEmitter<TorrentTabSelection>();

  @Input() selectedTab: TorrentTabSelection = undefined;

  transloco = inject(TranslocoService);
  grapql = inject(GraphQLService);
  errors = inject(ErrorsService);

  get selectedTabIndex(): number {
    return torrentTabNames.indexOf(this.selectedTab as TorrentTab) + 1;
  }

  selectTabIndex(index: number): void {
    this.selectedTab = torrentTabNames[index - 1];
    this.tabSelected.emit(this.selectedTab);
  }

  delete() {
    this.grapql
      .torrentDelete({ infoHashes: [this.torrentContent.infoHash] })
      .pipe(
        catchError((err: Error) => {
          this.errors.addError(`Error deleting torrent: ${err.message}`);
          return EMPTY;
        }),
      )
      .pipe(
        tap(() => {
          this.updated.emit(null);
        }),
      )
      .subscribe();
  }

  getAttribute(key: string, source?: string): string | undefined {
    return this.torrentContent.content?.attributes?.find(
      (a) => a.key === key && (source === undefined || a.source === source),
    )?.value;
  }

  getCollections(type: string): string[] | undefined {
    const collections = this.torrentContent.content?.collections
      ?.filter((a) => a.type === type)
      .map((a) => a.name);
    return collections?.length ? collections.sort() : undefined;
  }

  filesCount(): number | undefined {
    if (this.torrentContent.torrent.filesStatus === "single") {
      return 1;
    }
    return this.torrentContent.torrent.filesCount ?? undefined;
  }
}
