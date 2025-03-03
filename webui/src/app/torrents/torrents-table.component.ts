import {
  Component,
  EventEmitter,
  inject,
  Input,
  OnInit,
  Output,
} from "@angular/core";
import {
  animate,
  state,
  style,
  transition,
  trigger,
} from "@angular/animations";
import { SelectionModel } from "@angular/cdk/collections";
import { ActivatedRoute, Router } from "@angular/router";
import { FilesizePipe } from "../pipes/filesize.pipe";
import { TimeAgoPipe } from "../pipes/time-ago.pipe";
import * as generated from "../graphql/generated";
import { BreakpointsService } from "../layout/breakpoints.service";
import { AppModule } from "../app.module";
import { TorrentsSearchDatasource } from "./torrents-search.datasource";
import { contentTypeInfo } from "./content-types";
import { TorrentChipsComponent } from "./torrent-chips.component";
import { TorrentContentComponent } from "./torrent-content.component";
import { TorrentsSearchController } from "./torrents-search.controller";

@Component({
  selector: "app-torrents-table",
  standalone: true,
  imports: [
    AppModule,
    FilesizePipe,
    TimeAgoPipe,
    TorrentChipsComponent,
    TorrentContentComponent,
  ],
  templateUrl: "./torrents-table.component.html",
  styleUrl: "./torrents-table.component.scss",
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
export class TorrentsTableComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  breakpoints = inject(BreakpointsService);

  contentTypeInfo = contentTypeInfo;

  @Input() dataSource: TorrentsSearchDatasource;
  @Input() controller: TorrentsSearchController;
  @Input() multiSelection: SelectionModel<string>;
  @Input() displayedColumns: readonly Column[] = allColumns;

  @Output() updated = new EventEmitter<string>();

  items = Array<generated.TorrentContent>();

  ngOnInit() {
    this.dataSource.items$.subscribe((items) => {
      this.items = items;
    });
  }

  /** Whether the number of selected elements matches the total number of rows. */
  isAllSelected() {
    return this.items.every((i) => this.multiSelection.isSelected(i.infoHash));
  }

  /** Selects all rows if they are not all selected; otherwise clear selection. */
  toggleAllRows() {
    if (this.isAllSelected()) {
      this.multiSelection.clear();
      return;
    }
    this.multiSelection.select(...this.items.map((i) => i.infoHash));
  }

  toggleSelectedTorrent(infoHash: string) {
    this.controller.update((ctrl) => ({
      ...ctrl,
      selectedTorrent:
        ctrl.selectedTorrent?.infoHash === infoHash
          ? undefined
          : {
              infoHash,
              tab: ctrl.selectedTorrent?.tab,
            },
    }));
  }

  /**
   * Workaround for untyped table cell definitions
   */
  item(item: generated.TorrentContent): generated.TorrentContent {
    return item;
  }
}

export const allColumns = [
  "select",
  "summary",
  "size",
  "publishedAt",
  "peers",
  "magnet",
] as const;

export const compactColumns = ["select", "summary", "size", "magnet"] as const;

export type Column = (typeof allColumns)[number];
